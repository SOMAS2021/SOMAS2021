export default function Sidebar() {
  return (
    <div className="absolute dark:bg-gray-800 bg-gray-100">
      <div className="flex flex-col sm:flex-row sm:justify-around">
        <div className="w-72">
          <nav className="mt-10 px-6 scroll">
            {[...range(1, 100)].map((i) => (
              <SidebarElement title={`File no ${i}`} link={`/${i}`} />
            ))}
          </nav>
        </div>
      </div>
    </div>
  );
}

function range(start: number, end: number) {
  return Array(end - start + 1).fill(0).map((_, idx) => start + idx)
}

export interface SidebarElementProps {
  title: string;
  link: string;
}

function SidebarElement(props: SidebarElementProps) {
  const { title, link } = props;
  console.log(window.location.pathname);
  const active = link === window.location.pathname;
  return (
    <a
      className={`rounded-full hover:text-gray-800 hover:bg-gray-200 flex items-center p-2 my-6 transition-colors dark:hover:text-white dark:hover:bg-gray-600 duration-200 ${
        active ? "text-gray-800 dark:text-gray-100 bg-gray-200 dark:bg-gray-600" : "text-gray-600 dark:text-gray-400"
      }`}
      href={link}
    >
      <span className="mx-4 text-lg font-normal">{title}</span>
      <span className="flex-grow text-right"></span>
    </a>
  );
}
