use chrono::{Duration as ChronoDuration, Utc};
use identity_service::domain::models::session::Session;
use identity_service::service::session_cache_service::{
    SessionCacheService, SessionCacheServiceTrait,
};
use std::time::Duration;
use tokio_util::sync::CancellationToken;
use uuid::Uuid;

fn create_session(session_id: Uuid, user_id: Uuid, expires_in: ChronoDuration) -> Session {
    let now = Utc::now();
    Session {
        session_id,
        user_id,
        ip_address: Some("127.0.0.1".to_string()),
        user_agent: Some("TestAgent".to_string()),
        created_at: now,
        expires_at: now + expires_in,
        last_accessed_at: now,
    }
}

#[tokio::test]
async fn test_session_cache_basic_ops() {
    let shutdown_token = CancellationToken::new();
    let cache = SessionCacheService::new(Duration::from_secs(60), shutdown_token.clone(), 100);

    let session_id = Uuid::now_v7();
    let user_id = Uuid::now_v7();
    let session = create_session(session_id, user_id, ChronoDuration::hours(1));

    // Set and Get
    cache.set(session_id.to_string(), session.clone()).await;
    let cached = cache.get(session_id.to_string()).await;
    assert!(cached.is_some());
    assert_eq!(cached.unwrap().session_id, session_id);

    // Update last accessed
    cache.update_last_accessed(session_id.to_string()).await;

    // Remove
    cache.remove(session_id.to_string()).await;
    let cached = cache.get(session_id.to_string()).await;
    assert!(cached.is_none());

    shutdown_token.cancel();
}

#[tokio::test]
async fn test_session_cache_expiration() {
    let shutdown_token = CancellationToken::new();
    // Short cleanup interval for test
    let cache = SessionCacheService::new(Duration::from_millis(100), shutdown_token.clone(), 100);

    let session_id = Uuid::now_v7();
    let user_id = Uuid::now_v7();
    // Expire very soon
    let session = create_session(session_id, user_id, ChronoDuration::milliseconds(50));

    cache.set(session_id.to_string(), session).await;

    // Should be there initially
    assert!(cache.get(session_id.to_string()).await.is_some());

    // Wait for expiration and cleanup
    tokio::time::sleep(Duration::from_millis(200)).await;

    // Should be gone
    assert!(cache.get(session_id.to_string()).await.is_none());

    shutdown_token.cancel();
}

#[tokio::test]
async fn test_session_cache_capacity_eviction() {
    let shutdown_token = CancellationToken::new();
    let max_capacity = 2;
    let cache = SessionCacheService::new(
        Duration::from_secs(60),
        shutdown_token.clone(),
        max_capacity,
    );

    let user_id = Uuid::now_v7();

    let id1 = Uuid::now_v7().to_string();
    let mut s1 = create_session(Uuid::now_v7(), user_id, ChronoDuration::hours(1));
    s1.last_accessed_at = Utc::now() - ChronoDuration::minutes(10);

    let id2 = Uuid::now_v7().to_string();
    let mut s2 = create_session(Uuid::now_v7(), user_id, ChronoDuration::hours(1));
    s2.last_accessed_at = Utc::now() - ChronoDuration::minutes(5);

    cache.set(id1.clone(), s1).await;
    cache.set(id2.clone(), s2).await;

    // Capacity reached
    let id3 = Uuid::now_v7().to_string();
    let s3 = create_session(Uuid::now_v7(), user_id, ChronoDuration::hours(1));

    cache.set(id3.clone(), s3).await;

    // id1 (oldest last_accessed) should be evicted
    assert!(cache.get(id1).await.is_none());
    assert!(cache.get(id2).await.is_some());
    assert!(cache.get(id3).await.is_some());

    shutdown_token.cancel();
}

#[tokio::test]
async fn test_remove_user_sessions() {
    let shutdown_token = CancellationToken::new();
    let cache = SessionCacheService::new(Duration::from_secs(60), shutdown_token.clone(), 100);

    let user1 = Uuid::now_v7();
    let user2 = Uuid::now_v7();

    let s1 = create_session(Uuid::now_v7(), user1, ChronoDuration::hours(1));
    let s2 = create_session(Uuid::now_v7(), user1, ChronoDuration::hours(1));
    let s3 = create_session(Uuid::now_v7(), user2, ChronoDuration::hours(1));

    cache.set(s1.session_id.to_string(), s1.clone()).await;
    cache.set(s2.session_id.to_string(), s2.clone()).await;
    cache.set(s3.session_id.to_string(), s3.clone()).await;

    cache.remove_user_sessions(user1).await;

    assert!(cache.get(s1.session_id.to_string()).await.is_none());
    assert!(cache.get(s2.session_id.to_string()).await.is_none());
    assert!(cache.get(s3.session_id.to_string()).await.is_some());

    shutdown_token.cancel();
}
