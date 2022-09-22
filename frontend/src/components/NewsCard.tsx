type NewsCardProps = {
  category: string;
  categoryLink: string;
  title: string;
  date: string;
  author: string;
  URL: string;
  img: string;
  imgSrcSet: string;
};

export default function NewsCard({article}:{article:NewsCardProps}) {
  return (
    <div className="flex items-center my-5">
      <span
        className="uppercase -rotate-180 border-l text-[#3cffca] trackin-wider text-[.75rem] md:text-[1rem] cursor-pointer"
        style={{
          writingMode: "vertical-lr",
          textOrientation: "sideways",
        }}
      >
        <a href={article.categoryLink} target="_blank" rel="noreferrer">
          {article.category}
        </a>
      </span>
      <div className="flex-col ml-3 max-w-[80%]">
        <h2 className="font-bold  md:text-[1.37rem] cursor-pointer hover:underline sm:hover:underline-offset-2 md:hover:underline-offset-4 hover:decoration-[#3cffca] sm:w-[90%]">
          <a href={article.URL} target="_blank" rel="noreferrer">
            {article.title}
          </a>
        </h2>
        <div className="flex sm:text-[0.85rem] md:text-[1.1rem]">
          <p className="text-[#3cffca]">{article.author}</p>
          <p className="text-[#bdbdbd] ml-2">{article.date}</p>
        </div>
      </div>
      <div className="ml-auto max-w-[75px] rounded-[3px] sm:min-w-[100px] md:max-w-[150px] cursor-pointer">
        <a href={article.URL} target="_blank" rel="noreferrer">
          <img
            decoding={"async"}
            className="aspect-square block min-w-full h-full cover rounded-sm"
            sizes="240px"
            src={article.img}
            srcSet={article.imgSrcSet}
            loading="lazy"
          />
        </a>
      </div>
    </div>
  );
}
