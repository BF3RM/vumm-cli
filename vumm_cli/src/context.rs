use std::time::Duration;
use tokio::time::timeout;

#[allow(dead_code)]
pub struct CommandContext {
    pub timeout: Duration,
}

#[allow(dead_code)]
impl CommandContext {
    pub fn new(timeout: Duration) -> Self {
        Self { timeout }
    }

    pub async fn run_with_timeout<T, F, Fut>(self, f: F) -> Result<T, String>
    where
        F: FnOnce() -> Fut,
        Fut: std::future::Future<Output = T>,
    {
        match timeout(self.timeout, f()).await {
            Ok(result) => Ok(result),
            Err(_) => Err(format!("Operation timed out after {:?}", self.timeout)),
        }
    }
}
