import React, { useEffect, useState } from 'react';
import { StatusBar } from 'expo-status-bar';
import { Image, FlatList, StyleSheet, Text, View } from 'react-native';
import axios from 'axios';

const defaultImageURL = "https://upload.wikimedia.org/wikipedia/commons/a/ac/Default_pfp.jpg"

const ValidateImageURL = (url) => {
    fetch(url).then((res) => {
        if(res.status != 404) return
        url = defaultImageURL
    }).catch((err) => {
        url = defaultImageURL
    })
}

const OrgListRow = (props) => {
    const [validImage, setValidImage] = useState(true)
    var url = props.item.pfp_url
    return (
        <View style={{flexDirection: 'row', flexWrap:'wrap', paddingBottom: 15}}>
            <Image style={{width: 50, height: 50, borderRadius: 25}} source={{uri: validImage ? props.item.pfp_url : defaultImageURL}} onError={()=>{setValidImage(false)}}/>
            <Text>{ValidateImageURL(props.item.pfp_url)}</Text>
            <Text style={{justifyContent: 'center', paddingLeft: 20, fontSize: 20, paddingTop: 15}}>{props.item.lastName}, {props.item.firstName}</Text>
        </View>
    )
}

export default function OrgListView() {
    const [data, setData] = useState([])
    const [loaded, setLoaded] = useState(false)
  
    const getUsers = async () => {
      try {
        const response = await axios.post('http://192.168.1.12:8080/query', {
          query: `query{
            users{
              id
              firstName
              lastName
              pfp_url
              institution
            }
          }`
        })
        setData(response.data.data)
        console.log("data", response.data.data)
        setLoaded(true)
      }
      catch (error) {
        console.log(error)
      }
    }
  
    useEffect(() => {
      getUsers()
    }, [])
  
    if (!loaded) {
      return (
        <View></View>
      )
    }
    return (
      <View style={styles.container}>
        <FlatList data = {data.users} keyExtractor={(item) => item.id} renderItem = {({item}) => <OrgListRow item={item}/>}/>
        <StatusBar style="auto" />
      </View>
    );
  }

  const styles = StyleSheet.create({
    container: {
      backgroundColor: '#fff',
      paddingVertical: 200,
      paddingHorizontal: 50,
      justifyContent: 'center',
    },
  });
  