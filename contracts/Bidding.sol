// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Bidding {
    struct EncryptedBid {
        address bidder;
        bytes encryptedBid;
    }

    EncryptedBid[] public bids;

    function submitBid(bytes calldata encryptedBid) external {
        bids.push(EncryptedBid(msg.sender, encryptedBid));
    }

    function getBids() external view returns (EncryptedBid[] memory) {
        return bids;
    }
}