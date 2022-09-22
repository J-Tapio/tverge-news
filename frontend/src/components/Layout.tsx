import { ReactNode } from "react";
// Components
import Banner from "./Banner";

export default function Layout({children}:{children: ReactNode}) {
  return (
    <div id="app">
      <div className="w-full h-full pl-5 pr-8 bg-[#000000] overflow-y-hidden">
        <Banner />
        <div className="flex sm:h-[calc(100vh-18%)] md:h-[calc(100vh-22%)] lg:h-[calc(100vh-25%)] justify-end overflow-y-scroll scroll-smooth scrollbar-none">
          <div className="w-[80%] mb-10">
            <div className="sticky top-0 bg-[#000] opacity-50 h-[10px]" />
              {children}
          </div>
        </div>
      </div>
    </div>
  );
}