import '@near-wallet-selector/modal-ui/styles.css';

import React from 'react';
import './App.css'

import { setupWalletSelector } from "@near-wallet-selector/core";
import { setupModal } from "@near-wallet-selector/modal-ui";
import { setupMeteorWallet } from "@near-wallet-selector/meteor-wallet";
import { setupLedger } from "@near-wallet-selector/ledger";
import { WalletSelectorProvider } from "@near-wallet-selector/react-hook";
import BlockchainDataInfo from "./BlockchainDataInfo";


const selector = await setupWalletSelector({
  network: "testnet",
  modules: [
    setupMeteorWallet(),
    setupLedger(),
  ],
});

const walletSelectorConfig = {
  network: "testnet", 
  createAccessKeyFor: "vladozzzwrq.testnet",
  modules: [
    setupMeteorWallet(),
    setupLedger()
  ],
}


const modal = setupModal(selector, {
  contractId: "vladozzzwrq.testnet"
});
modal.show();


export default function App({ Component }) {


  return (
    <WalletSelectorProvider config={walletSelectorConfig}>
      <h1>Hello World</h1>
      <BlockchainDataInfo /> 
    </WalletSelectorProvider>
  );
}