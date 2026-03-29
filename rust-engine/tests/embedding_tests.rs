#[cfg(test)]
mod tests {
    use vector_engine::EmbeddingEngine;

    #[test]
    fn test_embed_single() {
        let engine = EmbeddingEngine::new().expect("Failed to create engine");
        let result = engine.embed_single("Hello world");
        
        assert!(result.is_ok());
        let embedding = result.unwrap();
        assert_eq!(embedding.len(), 384); // AllMiniLML6V2 dimension
    }

    #[test]
    fn test_cosine_similarity() {
        let a = vec![1.0, 0.0, 0.0];
        let b = vec![1.0, 0.0, 0.0];
        let c = vec![0.0, 1.0, 0.0];

        let sim_identical = EmbeddingEngine::cosine_similarity(&a, &b);
        let sim_orthogonal = EmbeddingEngine::cosine_similarity(&a, &c);

        assert!((sim_identical - 1.0).abs() < 0.001);
        assert!(sim_orthogonal.abs() < 0.001);
    }

    #[test]
    fn test_embed_batch() {
        let engine = EmbeddingEngine::new().expect("Failed to create engine");
        let texts = vec![
            "First document".to_string(),
            "Second document".to_string(),
        ];
        
        let result = engine.embed_batch(texts);
        assert!(result.is_ok());
        
        let embeddings = result.unwrap();
        assert_eq!(embeddings.len(), 2);
    }
}
