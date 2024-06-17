# Graceful

Graceful makes it easy to handle graceful shutdown in your Go applications.

- Automatically handle OS exit signals
- Blocking function to synchronize multiple goroutines
- Designed to work equally well with contextual and non-contextual goroutines
- The cancellation reason is provided upon exit
- No additional dependencies

```go
g := graceful.New(context.Background())

// Add goroutine using context
g.Go(func(ctx context.Context) error {
    return queue.Start(ctx)
})

// Add goroutine without using context
g.Go(func(ctx context.Context) error {
    return server.ListenAndServe()
})

// Hook to clean up contextless goroutines
g.Stop(func(ctx context.Context) error {
    return server.Shutdown(context.Background())
})

// Log the reason for stopping
log.Error().
    AnErr("cause", grace.Wait()).
    Msg("service shutdown")
```
