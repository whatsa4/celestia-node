/*
Package header contains all services related to generating, requesting, syncing and storing
ExtendedHeaders.

An ExtendedHeader is similar to a regular block header from celestia-core, but "extended" by adding
the DataAvailabilityHeader (DAH), ValidatorSet and Commit in order to allow for light and full nodes
to reconstruct block data by referencing the DAH.

ExtendedHeaders are generated by bridge nodes that listen to the celestia-core network for blocks,
validate and extend the block data in order to generate the DAH, and request both the ValidatorSet
and Commit for that block height from the celestia-core network. Bridge nodes then package that data
up into an ExtendedHeader and broadcast it to the 'header-sub' GossipSub topic (HeaderSub). All nodes subscribed
to HeaderSub will receive and validate the ExtendedHeader, and store it, making it available to all
other dependent services (such as the DataAvailabilitySampler, or DASer) to access.

There are 5 main components in the header package:
	1. core.Listener listens for new blocks from the celestia-core network (run by bridge nodes only),
	   extends them, generates a new ExtendedHeader, and publishes it to the HeaderSub.
	2. p2p.Subscriber listens for new ExtendedHeaders from the Celestia Data Availability (DA) network (via
	   the HeaderSub)
	3. p2p.Exchange or core.Exchange request ExtendedHeaders from other celestia DA nodes (default for
	   full and light nodes) or from a celestia-core node connection (bridge nodes only)
	4. Syncer manages syncing of past and recent ExtendedHeaders from either the DA network or a celestia-core
       connection (bridge nodes only).
	5. Store manages storing ExtendedHeaders and making them available for access by other dependent services.

For bridge nodes, the general flow of the header Service is as follows:
	1. core.Listener listens for new blocks from the celestia-core connection
	2. core.Listener validates the block and generates the ExtendedHeader, simultaneously storing the
       extended block shares to disk
	3. core.Listener publishes the new ExtendedHeader to HeaderSub, notifying all subscribed peers of the new
	   ExtendedHeader
	4. Syncer is already subscribed to the HeaderSub, so it receives new ExtendedHeaders locally from
       the core.Listener and stores them to disk via Store.
	4a. If the celestia-core connection is started simultaneously with the bridge node, then the celestia-core
       connection will handle the syncing component, piping every synced block to the core.Listener
	4b. If the celestia-core connection is already synced to the network, the Syncer handles requesting past
		headers up to the network head from the celestia-core connection (using core.Exchange rather than p2p.Exchange).

For full and light nodes, the general flow of the header Service is as follows:
	1. Syncer listens for new ExtendedHeaders from HeaderSub
	2. If there is a gap between the local head of chain and the new, validated network head, the Syncer
	   kicks off a sync routine to request all ExtendedHeaders between local head and network head.
	3. While the Syncer requests headers between the local head and network head in batches, it appends them to the
	   subjective chain via Store with the last batched header as the new local head.
*/
package header
