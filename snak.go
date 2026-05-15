import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  Alert,
  ActivityIndicator,
  FlatList,
  Image,
  ScrollView,
  Dimensions
} from 'react-native';

import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import { VideoView, useVideoPlayer } from 'expo-video';

// ================= CONFIG =================

const TOKEN = "KFB2YupQfqQZj6kFDwO6NKaje07vd6DP";

const BASE_URL =
  "http://213.199.56.115/api/database/rows/table/";

const TABLE_EPISODIOS = "1191";

const Stack = createNativeStackNavigator();

const { width } = Dimensions.get('window');

// ================= LOGIN =================

function LoginScreen({ navigation }) {

  const [email, setEmail] = useState('');
  const [senha, setSenha] = useState('');
  const [loading, setLoading] = useState(false);

  const fazerLogin = async () => {

    if (!email || !senha) {
      return Alert.alert(
        "Erro",
        "Preencha todos os campos"
      );
    }

    setLoading(true);

    try {

      const res = await fetch(
        `${BASE_URL}1199/?user_field_names=true`,
        {
          headers: {
            Authorization: `Token ${TOKEN}`
          }
        }
      );

      const data = await res.json();

      const usuario =
        data.results.find(
          u =>
            u.Email?.toLowerCase() ===
            email.toLowerCase() &&
            u.Senha === senha
        );

      if (usuario) {

        navigation.replace('Home');

      } else {

        Alert.alert(
          "Erro",
          "Email ou senha incorretos"
        );
      }

    } catch (e) {

      Alert.alert(
        "Erro",
        "Falha na conexão"
      );
    }

    setLoading(false);
  };

  return (
    <SafeAreaView style={styles.container}>

      <View style={styles.header}>

        <Text style={styles.logo}>
          PHFLIX
        </Text>

        <Text style={styles.subtitle}>
          Filmes e Séries
        </Text>

      </View>

      <View>

        <TextInput
          style={styles.input}
          placeholder="Email"
          placeholderTextColor="#888"
          value={email}
          onChangeText={setEmail}
          keyboardType="email-address"
          autoCapitalize="none"
        />

        <TextInput
          style={styles.input}
          placeholder="Senha"
          placeholderTextColor="#888"
          value={senha}
          onChangeText={setSenha}
          secureTextEntry
        />

        <TouchableOpacity
          style={styles.botao}
          onPress={fazerLogin}
        >

          {loading ? (

            <ActivityIndicator color="#fff" />

          ) : (

            <Text style={styles.textoBotao}>
              ENTRAR
            </Text>
          )}

        </TouchableOpacity>

      </View>

    </SafeAreaView>
  );
}

// ================= DETAILS =================

function DetailsScreen({ route, navigation }) {

  const { conteudo } = route.params;

  const [episodios, setEpisodios] = useState([]);
  const [loadingEp, setLoadingEp] = useState(false);

  // ================= AJUSTE =================
  // aceita Serie, Série etc
  // ==========================================

  const isSerie =
    conteudo.Tipo
      ?.toString()
      .toLowerCase()
      .includes('s');

  useEffect(() => {

    if (isSerie) {
      carregarEpisodios();
    }

  }, []);

  const carregarEpisodios = async () => {

    setLoadingEp(true);

    try {

      const url =
        `${BASE_URL}${TABLE_EPISODIOS}/?user_field_names=true&size=500`;

      const res = await fetch(url, {
        headers: {
          Authorization: `Token ${TOKEN}`
        }
      });

      const data = await res.json();

      const lista =
        data.results || [];
       
      console.log(JSON.stringify(lista[0], null, 2));

      console.log("TODOS EPISODIOS:", lista);

      // ================= FILTRO =================

      const filtrados = lista.filter(ep => {

  const nomeSerie =
    String(conteudo.Nome || '')
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .trim()
      .toLowerCase();

  const nomeEp =
    String(ep.Nome || '')
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .trim()
      .toLowerCase();

  console.log("SERIE:", nomeSerie);
  console.log("EP:", nomeEp);

  return nomeEp.includes(nomeSerie);
});

      // ================= ORDENAÇÃO =================

      filtrados.sort((a, b) => {

        const tempA =
          Number(a.Temporada || 0);

        const tempB =
          Number(b.Temporada || 0);

        if (tempA !== tempB) {
          return tempA - tempB;
        }

        return (
          Number(a.Episodio || 0) -
          Number(b.Episodio || 0)
        );
      });

      console.log("EPISODIOS FILTRADOS:", filtrados);

      setEpisodios(filtrados);

    } catch (e) {

      console.log(e);

      Alert.alert(
        "Erro",
        "Não foi possível carregar episódios"
      );
    }

    setLoadingEp(false);
  };

  const playVideo = (link) => {

    if (!link) {

      return Alert.alert(
        "Erro",
        "Link não disponível"
      );
    }

    navigation.navigate(
      'Player',
      { link }
    );
  };

  // ================= AGRUPAR TEMPORADAS =================

  const temporadas = {};

  episodios.forEach(ep => {

    const temp =
      ep.Temporada || 1;

    if (!temporadas[temp]) {
      temporadas[temp] = [];
    }

    temporadas[temp].push(ep);
  });

  return (
    <SafeAreaView style={styles.detailsContainer}>

      <ScrollView>

        <Image
          source={{ uri: conteudo.Capa }}
          style={styles.detailPoster}
          resizeMode="cover"
        />

        <View style={styles.detailsContent}>

          <Text style={styles.detailTitle}>
            {conteudo.Nome}
          </Text>

          {conteudo.Categoria && (

            <Text style={styles.genre}>
              {conteudo.Categoria}
            </Text>
          )}

          {/* ================= SERIES ================= */}

          {isSerie && (

            <View>

              <Text style={styles.seasonTitle}>
                Temporadas
              </Text>

              {loadingEp ? (

                <ActivityIndicator
                  color="#E50914"
                  size="large"
                />

              ) : episodios.length > 0 ? (

                Object.keys(temporadas).map(temp => (

                  <View key={temp}>

                    <Text
                      style={{
                        color: '#E50914',
                        fontSize: 22,
                        fontWeight: 'bold',
                        marginTop: 20,
                        marginBottom: 10
                      }}
                    >
                      TEMPORADA {temp}
                    </Text>

                    {temporadas[temp].map((ep, i) => (

                      <TouchableOpacity
                        key={i}
                        style={styles.episodeItem}
                        onPress={() =>
                          playVideo(ep.Link)
                        }
                      >

                        <View>

                          <Text style={styles.episodeText}>
                            EP {ep.Episodio}
                          </Text>

                          <Text
                            style={{
                              color: '#999',
                              marginTop: 4
                            }}
                          >
                            {ep.Nome}
                          </Text>

                        </View>

                        <Text
                          style={{
                            color: '#E50914',
                            fontSize: 22
                          }}
                        >
                          ▶
                        </Text>

                      </TouchableOpacity>
                    ))}

                  </View>
                ))

              ) : (

                <Text
                  style={{
                    color: '#ff4444',
                    textAlign: 'center',
                    marginTop: 20
                  }}
                >
                  Nenhum episódio encontrado
                </Text>
              )}

            </View>
          )}

          {/* ================= FILMES ================= */}

          {!isSerie && conteudo.Link && (

            <TouchableOpacity
              style={styles.playButton}
              onPress={() =>
                playVideo(conteudo.Link)
              }
            >

              <Text style={styles.playButtonText}>
                ▶ ASSISTIR AGORA
              </Text>

            </TouchableOpacity>
          )}

        </View>

      </ScrollView>

    </SafeAreaView>
  );
}

// ================= PLAYER =================

function PlayerScreen({ route, navigation }) {

  const { link } =
    route.params || {};

  const player =
    useVideoPlayer(link, (p) => {
      p.loop = false;
    });

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: '#000'
      }}
    >

      <TouchableOpacity
        onPress={() =>
          navigation.goBack()
        }
        style={styles.backButton}
      >

        <Text
          style={{
            color: '#fff',
            fontSize: 18
          }}
        >
          ← Voltar
        </Text>

      </TouchableOpacity>

      <VideoView
        player={player}
        style={{ flex: 1 }}
        allowsFullscreen
        allowsPictureInPicture
        contentFit="contain"
      />

    </SafeAreaView>
  );
}

// ================= HOME =================

function HomeScreen({ navigation }) {

  const [banners, setBanners] = useState([]);
  const [conteudos, setConteudos] = useState([]);
  const [categorias, setCategorias] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    carregarTudo();
  }, []);

  const carregarTudo = async () => {

    try {

      const [b, c, cat] =
        await Promise.all([

          fetch(
            `${BASE_URL}1189/?user_field_names=true`,
            {
              headers: {
                Authorization: `Token ${TOKEN}`
              }
            }
          ),

          fetch(
            `${BASE_URL}1186/?user_field_names=true&size=100`,
            {
              headers: {
                Authorization: `Token ${TOKEN}`
              }
            }
          ),

          fetch(
            `${BASE_URL}1190/?user_field_names=true`,
            {
              headers: {
                Authorization: `Token ${TOKEN}`
              }
            }
          )
        ]);

      setBanners(
        (await b.json()).results || []
      );

      setConteudos(
        (await c.json()).results || []
      );

      setCategorias(
        (await cat.json()).results || []
      );

    } catch (e) {

      Alert.alert(
        "Erro",
        "Falha ao carregar dados"
      );
    }

    setLoading(false);
  };

  const abrirDetalhes = (item) => {

    navigation.navigate(
      'Details',
      { conteudo: item }
    );
  };

  const renderItem = ({ item }) => (

    <TouchableOpacity
      style={styles.card}
      onPress={() =>
        abrirDetalhes(item)
      }
    >

      <Image
        source={{ uri: item.Capa }}
        style={styles.capa}
        resizeMode="cover"
      />

      <Text
        style={styles.tituloCard}
        numberOfLines={2}
      >
        {item.Nome}
      </Text>

    </TouchableOpacity>
  );

  if (loading) {

    return (
      <SafeAreaView style={styles.center}>

        <ActivityIndicator
          size="large"
          color="#E50914"
        />

      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.homeContainer}>

      <ScrollView>

        <View style={styles.bannerContainer}>

          <FlatList
            data={banners}
            horizontal
            pagingEnabled
            showsHorizontalScrollIndicator={false}
            keyExtractor={(item, index) =>
              index.toString()
            }
            renderItem={({ item }) => (

              <Image
                source={{ uri: item.Imagem }}
                style={styles.banner}
                resizeMode="cover"
              />
            )}
          />

          <View style={styles.logoOverlay}>

            <Text style={styles.logoHome}>
              PHFLIX
            </Text>

          </View>

        </View>

        {categorias.map(cat => {

          const filtered =
            conteudos.filter(
              con =>
                con.Categoria &&
                con.Categoria.includes(cat.Nome)
            );

          if (filtered.length === 0)
            return null;

          return (

            <View key={cat.id}>

              <Text style={styles.secaoTitulo}>
                {cat.Nome}
              </Text>

              <FlatList
                data={filtered}
                renderItem={renderItem}
                horizontal
                keyExtractor={(item, index) =>
                  index.toString()
                }
                showsHorizontalScrollIndicator={false}
                contentContainerStyle={styles.lista}
              />

            </View>
          );
        })}

      </ScrollView>

    </SafeAreaView>
  );
}

// ================= APP =================

export default function App() {

  return (
    <SafeAreaProvider>

      <NavigationContainer>

        <Stack.Navigator
          screenOptions={{
            headerShown: false
          }}
        >

          <Stack.Screen
            name="Login"
            component={LoginScreen}
          />

          <Stack.Screen
            name="Home"
            component={HomeScreen}
          />

          <Stack.Screen
            name="Details"
            component={DetailsScreen}
          />

          <Stack.Screen
            name="Player"
            component={PlayerScreen}
          />

        </Stack.Navigator>

      </NavigationContainer>

    </SafeAreaProvider>
  );
}

// ================= STYLES =================

const styles = StyleSheet.create({

  container: {
    flex: 1,
    backgroundColor: '#000',
    justifyContent: 'center',
    padding: 20
  },

  header: {
    alignItems: 'center',
    marginBottom: 50
  },

  logo: {
    fontSize: 52,
    fontWeight: 'bold',
    color: '#E50914',
    letterSpacing: 4
  },

  subtitle: {
    fontSize: 18,
    color: '#fff',
    marginTop: 8
  },

  input: {
    backgroundColor: '#1F1F1F',
    color: '#fff',
    padding: 18,
    borderRadius: 8,
    marginBottom: 16,
    fontSize: 16
  },

  botao: {
    backgroundColor: '#E50914',
    padding: 18,
    borderRadius: 8,
    alignItems: 'center'
  },

  textoBotao: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold'
  },

  homeContainer: {
    flex: 1,
    backgroundColor: '#000'
  },

  bannerContainer: {
    position: 'relative'
  },

  banner: {
    width: width,
    height: 230
  },

  logoOverlay: {
    position: 'absolute',
    top: 40,
    left: 20,
    zIndex: 10
  },

  logoHome: {
    fontSize: 45,
    fontWeight: 'bold',
    color: '#E50914',
    letterSpacing: 3
  },

  secaoTitulo: {
    color: '#fff',
    fontSize: 22,
    fontWeight: 'bold',
    marginLeft: 15,
    marginTop: 25,
    marginBottom: 10
  },

  lista: {
    paddingLeft: 15
  },

  card: {
    marginRight: 12,
    width: 150
  },

  capa: {
    width: 150,
    height: 210,
    borderRadius: 8
  },

  tituloCard: {
    color: '#ddd',
    marginTop: 8,
    fontSize: 13.5,
    textAlign: 'center'
  },

  detailsContainer: {
    flex: 1,
    backgroundColor: '#000'
  },

  detailPoster: {
    width: '100%',
    height: 350
  },

  detailsContent: {
    padding: 15
  },

  detailTitle: {
    color: '#fff',
    fontSize: 26,
    fontWeight: 'bold',
    marginBottom: 8
  },

  genre: {
    color: '#E50914',
    fontSize: 16,
    marginBottom: 15
  },

  seasonTitle: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
    marginVertical: 15
  },

  episodeItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    backgroundColor: '#1F1F1F',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10
  },

  episodeText: {
    color: '#fff',
    fontSize: 16,
    flex: 1
  },

  playButton: {
    backgroundColor: '#E50914',
    padding: 18,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 20
  },

  playButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold'
  },

  center: {
    flex: 1,
    backgroundColor: '#000',
    justifyContent: 'center',
    alignItems: 'center'
  },

  backButton: {
    position: 'absolute',
    top: 50,
    left: 20,
    zIndex: 10,
    backgroundColor: 'rgba(0,0,0,0.7)',
    paddingHorizontal: 15,
    paddingVertical: 8,
    borderRadius: 8
  }

});
