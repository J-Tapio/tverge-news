import TvergeArticle from "./@types";
import { useState } from "react";
// Wails
import { EventsOn } from "../wailsjs/runtime";
import { Latest } from "../wailsjs/go/main/App";
// Components
import Layout from "./components/Layout";
import News from "./components/News";
import Loader from "./components/Loader";

function App() {
  const [newsLoaded, setNewsLoaded] = useState<boolean>(false);
  const [news, setNews] = useState<TvergeArticle[]>([]);

  async function getNews() {
    try {
      let latestNews = await Latest();
      setNews(latestNews);
      setNewsLoaded(true);
    } catch (error) {
      console.error(error);
    }
  }

  EventsOn("news", getNews);

  return (
    <>
      {!newsLoaded && <Loader />}
      {
        newsLoaded && (
        <Layout>
          <News news={news} />
        </Layout>
        )
      }
    </>
  );
}

export default App;
