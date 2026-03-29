use anyhow::Result;
use qdrant_client::prelude::*;
use qdrant_client::qdrant::{PointStruct, SearchPoints};
use serde_json::Value;
use std::collections::HashMap;

pub struct SearchEngine {
    client: QdrantClient,
}

impl SearchEngine {
    pub async fn new(url: &str) -> Result<Self> {
        let client = QdrantClient::from_url(url).build()?;
        Ok(Self { client })
    }

    pub async fn search(
        &self,
        collection: &str,
        query_vector: Vec<f32>,
        limit: u64,
    ) -> Result<Vec<ScoredPoint>> {
        let results = self
            .client
            .search_points(&SearchPoints {
                collection_name: collection.to_string(),
                vector: query_vector,
                limit,
                with_payload: Some(true.into()),
                ..Default::default()
            })
            .await?;

        Ok(results.result)
    }

    pub async fn upsert_points(
        &self,
        collection: &str,
        points: Vec<PointStruct>,
    ) -> Result<()> {
        self.client
            .upsert_points_blocking(collection, None, points, None)
            .await?;
        Ok(())
    }
}
