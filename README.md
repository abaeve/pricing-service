# pricing-service

The pricing service currently has 2 responsibilities:

 * List to what the pricing-fetcher is throwing onto relevant topics
   * Storing the current state of things, all available orders on the market and no more
 * Generating and making available stats on the above data

No history is stored here although there are plans to publish the changes which happened since the last fetch from pricing-fetcher

## Immediate roadmap

 * Publish an endpoint to control monitored regions

## Thoughts

 * This might be seen as the steward of pricing-fetcher
   * It could be made more aware of it's existence and pricing-fetcher could become a more robust service
     with endpoints to control what IT monitors on it's poll executions
   * This would allow for easier fine tuning of network and cpu usage
