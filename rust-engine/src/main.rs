use anyhow::Result;
use fastembed::{EmbeddingModel, InitOptions, TextEmbedding};
use qdrant_client::prelude::*;
use qdrant_client::qdrant::{CreateCollection, Distance, VectorParams, VectorsConfig};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{transport::Server, Request, Response, Status};

pub mod vector_service {
    tonic::include_proto!("vectorservice");
}

use vector_service::vector_service_server::{VectorService, VectorServiceServer};
use vector_service::{EmbedRequest, EmbedResponse, SearchRequest, SearchResponse};

#[derive(Clone)]
struct VectorEngineService {
    model: Arc<RwLock<TextEmbedding>>,
    qdrant: Arc<QdrantClient>,
}

impl VectorEngineService {
    async fn new() -> Result<Self> {
        let model = TextEmbedding::try_new(InitOptions {
            model_name: EmbeddingModel::AllMiniLML6V2,
            show_download_progress: true,
            ..Default::default()
        })?;

        let qdrant = QdrantClient::from_url("http://localhost:6334").build()?;

        Ok(Self {
            model: Arc::new(RwLock::new(model)),
            qdrant: Arc::new(qdrant),
        })
    }

    async fn ensure_collection(&self, collection_name: &str) -> Result<()> {
        let collections = self.qdrant.list_collections().await?;
        
        if !collections.collections.iter().any(|c| c.name == collection_name) {
            self.qdrant
                .create_collection(&CreateCollection {
                    collection_name: collection_name.to_string(),
                    vectors_config: Some(VectorsConfig {
                        config: Some(qdrant_client::qdrant::vectors_config::Config::Params(
                            VectorParams {
                                size: 384,
                                distance: Distance::Cosine.into(),
                                ..Default::default()
                            },
                        )),
                    }),
                    ..Default::default()
                })
                .await?;
        }
        Ok(())
    }
}

#[tonic::async_trait]
impl VectorService for VectorEngineService {
    async fn embed(
        &self,
        request: Request<EmbedRequest>,
    ) -> Result<Response<EmbedResponse>, Status> {
        let req = request.into_inner();
        let model = self.model.read().await;
        
        let embeddings = model
            .embed(vec![req.text], None)
            .map_err(|e| Status::internal(format!("Embedding failed: {}", e)))?;

        let vector = embeddings[0].to_vec();
        
        Ok(Response::new(EmbedResponse { vector }))
    }

    async fn search(
        &self,
        request: Request<SearchRequest>,
    ) -> Result<Response<SearchResponse>, Status> {
        let req = request.into_inner();
        
        self.ensure_collection(&req.collection)
            .await
            .map_err(|e| Status::internal(format!("Collection error: {}", e)))?;

        let results = self
            .qdrant
            .search_points(&qdrant_client::qdrant::SearchPoints {
                collection_name: req.collection,
                vector: req.query_vector,
                limit: req.limit as u64,
                with_payload: Some(true.into()),
                ..Default::default()
            })
            .await
            .map_err(|e| Status::internal(format!("Search failed: {}", e)))?;

        let matches = results
            .result
            .into_iter()
            .map(|point| vector_service::SearchMatch {
                id: point.id.map(|id| id.to_string()).unwrap_or_default(),
                score: point.score,
                payload: serde_json::to_string(&point.payload).unwrap_or_default(),
            })
            .collect();

        Ok(Response::new(SearchResponse { matches }))
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let addr = "[::1]:50051".parse()?;
    let service = VectorEngineService::new().await?;

    println!("Vector Engine gRPC server listening on {}", addr);

    Server::builder()
        .add_service(VectorServiceServer::new(service))
        .serve(addr)
        .await?;

    Ok(())
}
