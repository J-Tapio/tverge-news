export namespace main {
	
	export class TvergeArticle {
	    category: string;
	    categoryLink: string;
	    title: string;
	    date: string;
	    author: string;
	    URL: string;
	    img: string;
	    imgSrcSet: string;
	
	    static createFrom(source: any = {}) {
	        return new TvergeArticle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.categoryLink = source["categoryLink"];
	        this.title = source["title"];
	        this.date = source["date"];
	        this.author = source["author"];
	        this.URL = source["URL"];
	        this.img = source["img"];
	        this.imgSrcSet = source["imgSrcSet"];
	    }
	}

}

