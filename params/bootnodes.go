// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// GNC Foundation Go Bootnodes
	"enode://7ec759c185382e169e3fbd4718ca9907f26897109a5401aadb6128c26bce7bc9bc26613d202e67db6019137bec412d87a1f1411877cd5e165a544515c57203c9@47.57.116.216:30303",
	"enode://b09b539b1bf9b0bb4545fe89b969e440361a6c790809f8769d1bc91e4d631462e21be17db6db9f1df67f0380cc23ff39e32843496ca501e7d5022ac378ec7484@47.57.115.222:30303",
	"enode://d01987c09dc1149c7da115c8e9335b531716510fce9a0dec74d226a8a6d582cbbf8e455400e5c446359bee7e92341a0a62f69004474b20e9fcac302bed478c32@47.75.203.235:30303",
	"enode://10b45f2f3d3f27a77d19c74caabb24027814beea634ca540828ced46cf726b0808c59a3aafe7d76955fe55b5d4a6fb749b3611d27529fb0802a7ff0bccbac505@47.115.113.25:30303",
	"enode://8360aa1727dc16d9e977cae050c2966d15cabd7bc5167089fb52c685522c37be822c86f62f1f7b6370b8cca16a321b0c7eaa4dab7386dec58a65461da48c8a5c@8.208.98.131:30303",
	"enode://9012b9d0d7eed8d7ae3bdeca540ad2ca38e59ae65645adb8beb1b80f695cc69c7d34e9ef20d86475914200222a81b7f23eb3e4208e2ca2b55ed0395b1c397110@47.253.4.226:30303",
	"enode://3a5ad88cca847b07d8b2a0110997e74eab115c545db55a48aef8afd6a334f979d77dbee4e5519345aa606e22a7a2777229c3f355345dd6437a4307037e308d08@47.89.246.33:30303",
	"enode://5d6f43ef85c4822fa1dd66d6b3b260059773fe5bbe40d890b83fdd06e9daa5a07041577981c342419444b9965fa33720ac2552f0f0705153c64c83a3acd1587f@47.74.53.218:30303",

}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	"enode://30b7ab30a01c124a6cceca36863ece12c4f5fa68e3ba9b0b51407ccc002eeed3b3102d20a88f1c1d3c3154e2449317b8ef95090e77b312d5cc39354f86d5d606@52.176.7.10:30303",    // US-Azure geth
	"enode://865a63255b3bb68023b6bffd5095118fcc13e79dcf014fe4e47e065c350c7cc72af2e53eff895f11ba1bbb6a2b33271c1116ee870f266618eadfc2e78aa7349c@52.176.100.77:30303",  // US-Azure parity
	"enode://6332792c4a00e3e4ee0926ed89e0d27ef985424d97b6a45bf0f23e51f0dcb5e66b875777506458aea7af6f9e4ffb69f43f3778ee73c81ed9d34c51c4b16b0b0f@52.232.243.152:30303", // Parity
	"enode://94c15d1b9e2fe7ce56e458b9a3b672ef11894ddedd0c6f247e0f1d3487f52b66208fb4aeb8179fce6e3a749ea93ed147c37976d67af557508d199d9594c35f09@192.81.208.223:30303", // @gpip
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
	"enode://343149e4feefa15d882d9fe4ac7d88f885bd05ebb735e547f12e12080a9fa07c8014ca6fd7f373123488102fe5e34111f8509cf0b7de3f5b44339c9f25e87cb8@52.3.158.184:30303",  // INFURA
	"enode://b6b28890b006743680c52e64e0d16db57f28124885595fa03a562be1d2bf0f3a1da297d56b13da25fb992888fd556d4c1a27b1f39d531bde7de1921c90061cc6@159.89.28.211:30303", // AKASHA
}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Görli test network.
var GoerliBootnodes = []string{
	// Upstream bootnodes
	"enode://011f758e6552d105183b1761c5e2dea0111bc20fd5f6422bc7f91e0fabbec9a6595caf6239b37feb773dddd3f87240d99d859431891e4a642cf2a0a9e6cbb98a@51.141.78.53:30303",
	"enode://176b9417f511d05b6b2cf3e34b756cf0a7096b3094572a8f6ef4cdcb9d1f9d00683bf0f83347eebdf3b81c3521c2332086d9592802230bf528eaf606a1d9677b@13.93.54.137:30303",
	"enode://46add44b9f13965f7b9875ac6b85f016f341012d84f975377573800a863526f4da19ae2c620ec73d11591fa9510e992ecc03ad0751f53cc02f7c7ed6d55c7291@94.237.54.114:30313",
	"enode://c1f8b7c2ac4453271fa07d8e9ecf9a2e8285aa0bd0c07df0131f47153306b0736fd3db8924e7a9bf0bed6b1d8d4f87362a71b033dc7c64547728d953e43e59b2@52.64.155.147:30303",
	"enode://f4a9c6ee28586009fb5a96c8af13a58ed6d8315a9eee4772212c1d4d9cebe5a8b8a78ea4434f318726317d04a3f531a1ef0420cf9752605a562cfe858c46e263@213.186.16.82:30303",

	// Ethereum Foundation bootnode
	"enode://573b6607cd59f241e30e4c4943fd50e99e2b6f42f9bd5ca111659d309c06741247f4f1e93843ad3e8c8c18b6e2d94c161b7ef67479b3938780a97134b618b5ce@52.56.136.200:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://7ec759c185382e169e3fbd4718ca9907f26897109a5401aadb6128c26bce7bc9bc26613d202e67db6019137bec412d87a1f1411877cd5e165a544515c57203c9@47.57.116.216:30303",
	"enode://b09b539b1bf9b0bb4545fe89b969e440361a6c790809f8769d1bc91e4d631462e21be17db6db9f1df67f0380cc23ff39e32843496ca501e7d5022ac378ec7484@47.57.115.222:30303",
	"enode://d01987c09dc1149c7da115c8e9335b531716510fce9a0dec74d226a8a6d582cbbf8e455400e5c446359bee7e92341a0a62f69004474b20e9fcac302bed478c32@47.75.203.235:30303",
	"enode://10b45f2f3d3f27a77d19c74caabb24027814beea634ca540828ced46cf726b0808c59a3aafe7d76955fe55b5d4a6fb749b3611d27529fb0802a7ff0bccbac505@47.115.113.25:30303",
	"enode://8360aa1727dc16d9e977cae050c2966d15cabd7bc5167089fb52c685522c37be822c86f62f1f7b6370b8cca16a321b0c7eaa4dab7386dec58a65461da48c8a5c@8.208.98.131:30303",
	"enode://9012b9d0d7eed8d7ae3bdeca540ad2ca38e59ae65645adb8beb1b80f695cc69c7d34e9ef20d86475914200222a81b7f23eb3e4208e2ca2b55ed0395b1c397110@47.253.4.226:30303",
	"enode://3a5ad88cca847b07d8b2a0110997e74eab115c545db55a48aef8afd6a334f979d77dbee4e5519345aa606e22a7a2777229c3f355345dd6437a4307037e308d08@47.89.246.33:30303",
	"enode://5d6f43ef85c4822fa1dd66d6b3b260059773fe5bbe40d890b83fdd06e9daa5a07041577981c342419444b9965fa33720ac2552f0f0705153c64c83a3acd1587f@47.74.53.218:30303",
}
