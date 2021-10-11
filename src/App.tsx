import React, { useEffect, useState } from "react";
import { BiDownload, BiHdd, BiTimer } from "react-icons/bi";
import { Link, Route, Switch } from "react-router-dom";
import useFetch, { CachePolicies } from "use-http";
import create from "zustand";
import ReplayDetails, { Replay } from "./ReplayDetails";

export const searchStore = create<
  { value: string; setValue: (value: string) => void }
>((set) => {
  return {
    value: "",
    setValue: (value: string) =>
      set((state) => {
        return { ...state, value };
      }),
  };
});

function ReplayItem({ replay }: { replay: Replay }) {
  const time = new Date(Date.parse(replay.recording_time));
  const measuredTime = new Date(0);
  measuredTime.setSeconds(replay.duration);
  const parts = replay.path.split("/");

  return (
    <Link
      to={`/replay/${replay.id}`}
      className="p-2 border border-blue-300 rounded-sm hover:bg-blue-200 w-full flex flex-col"
    >
      <div className="flex flex-row text-lg mb-2">
        {time.toLocaleString()}
        <span className="ml-auto">{parts[parts.length - 1]}</span>
      </div>
      <div className="flex flex-row">
        <a
          className="p-2 border border-green-500 bg-green-200 hover:bg-green-300 text-green-700 rounded-sm shadow-sm text-sm my-auto"
          onClick={(e) => e.stopPropagation()}
          href={`/api/replay/${replay.id}/download`}
        >
          Download
          <BiDownload className="inline-flex w-4 h-4 ml-1" />
        </a>
        <div className="ml-auto flex flex-col">
          <div className="flex flex-row">
            <span className="text-gray-800 mr-2">
              {measuredTime.toISOString().substr(11, 8)}
            </span>
            <BiTimer className="inline-flex w-5 h-5 text-gray-500" />
            {" "}
          </div>
          <div className="flex flex-row">
            <span className="ml-auto text-gray-800 mr-2">
              {Math.round(replay.size / 1024 / 1024)} mb
            </span>
            <BiHdd className="inline-flex w-5 h-5 text-gray-500" />
            {" "}
          </div>
        </div>
      </div>
    </Link>
  );
}

function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(
    () => {
      const handler = setTimeout(() => {
        setDebouncedValue(value);
      }, delay);

      return () => {
        clearTimeout(handler);
      };
    },
    [value, delay],
  );

  return debouncedValue;
}

function ReplayList() {
  const [search, setSearch] = searchStore(
    (state) => [state.value, state.setValue],
  );
  const debouncedSearch = useDebounce(search, 500);
  const { data } = useFetch<Array<Replay>>(
    `/api/replay${
      debouncedSearch !== ""
        ? `?filter=${encodeURIComponent(debouncedSearch)}`
        : ""
    }`,
    {
      cachePolicy: CachePolicies.NO_CACHE,
    },
    [debouncedSearch],
  );

  if (!data) return <></>;
  return (
    <div className="m-auto md:w-1/4 w-full">
      <div className="p-2">
        <input
          className="border border-gray-300 rounded-sm form-input w-full h-8"
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>
      <div
        className="border border-gray-300 bg-gray-200 rounded-sm p-2 shadow-sm flex flex-col items-center gap-2 md:m-4"
      >
        {data.map((it) => <ReplayItem key={it.id} replay={it} />)}
      </div>
    </div>
  );
}

function App() {
  return (
    <div
      className="flex h-screen w-screen"
    >
      <Switch>
        <Route exact path="/" component={ReplayList} />
        <Route
          exact
          path="/replay/:replayId"
          render={({ match: { params: { replayId } } }) => {
            return <ReplayDetails replayId={parseInt(replayId)} />;
          }}
        />
      </Switch>
    </div>
  );
}

export default App;
