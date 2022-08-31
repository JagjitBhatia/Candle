import React, { useEffect, useState } from 'react';
import { StatusBar } from 'expo-status-bar';
import { FlatList, StyleSheet, Text, View } from 'react-native';
import axios from 'axios';

import OrgListView from './components/views/OrgListView';

export default function App() {
  return (
    <OrgListView/>
  )
}

