import { useEffect, useState } from 'react';
import { Latest } from '../../wailsjs/go/main/App'
// Components
import NewsCard from './NewsCard'

interface TvergeArticle {
  category: string;
  categoryLink: string;
  title: string;
  date: string;
  author: string;
  URL: string;
  img: string;
  imgSrcSet: string;
}

export default function News() {
  const [news, setNews] = useState<TvergeArticle[]>([]);

  async function getNews() {
    try {
      let latestNews = await Latest()
      console.log(latestNews)
      setNews(latestNews)
    } catch (error) {
      console.error(error);
    }
  }

  useEffect(() => {
    getNews()
  }, [])

  return (
    <>
    {news.length > 0 && news.map((article) => {
      return <NewsCard key={article.URL} article={article}/>
    })}
    </>
  )
}