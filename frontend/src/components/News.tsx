import TvergeArticle from '../@types';
// Components
import NewsCard from './NewsCard'


type NewsProps = {
  news: TvergeArticle[];
};

export default function News({news}: NewsProps) {
  return (
    <div className="pt-6">
      {news.length > 0 &&
        news.map((article) => {
          return <NewsCard key={article.URL} article={article} />;
        })}
    </div>
  );
}