import { ReactNode } from "react";
// Components
import Banner from "./Banner";
// The Transition component will render a div by default
import { Transition } from "@headlessui/react";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div id="app">
      <Transition
        show={true}
        enter="transition-opacity duration-150"
        enterFrom="opacity-0"
        enterTo="opacity-100"
        className="w-full h-full pl-5 pr-8 overflow-y-hidden"
      >
        <Banner />
        <div className="flex sm:h-[calc(100vh-18%)] md:h-[calc(100vh-22%)] lg:h-[calc(100vh-25%)] justify-end overflow-y-scroll scrollbar-none">
          <div className="w-[80%] curtain">
            {children}
          </div>
        </div>
      </Transition>
    </div>
  );
}
