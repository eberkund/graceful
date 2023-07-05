# Graceful

Graceful makes it easy to handle graceful shutdown in your Go applications.

- Automatically handle OS exit signals
- Blocking function to syncronize multiple goroutines
- Designed to work equally well with contextual and non-contextual goroutines
- The cancellation reason is provided upon exit

```go
grace, ctx := graceful.New(context.Background())

// Add goroutine using context
grace.Go(func() error {
    return queue.Start(ctx)
})

// Add goroutine without using context
grace.Go(func() error {
    return server.ListenAndServe()
})

// Hook to cleanup contextless goroutines
grace.Cleanup(func() error {
    return server.Shutdown(context.Background())
})

// Log the reason for stopping
log.Error().
    AnErr("cause", grace.Wait()).
    Msg("service shutdown")
```
