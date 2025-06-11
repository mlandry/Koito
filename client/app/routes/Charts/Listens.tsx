import ChartLayout from "./ChartLayout";
import { Link, useLoaderData, type LoaderFunctionArgs } from "react-router";
import { type Album, type Listen, type PaginatedResponse } from "api/api";
import { timeSince } from "~/utils/utils";
import ArtistLinks from "~/components/ArtistLinks";

export async function clientLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url);
  const page = url.searchParams.get("page") || "0";
  url.searchParams.set('page', page)

  const res = await fetch(
    `/apis/web/v1/listens?${url.searchParams.toString()}`
  );
  if (!res.ok) {
    throw new Response("Failed to load top tracks", { status: 500 });
  }

  const listens: PaginatedResponse<Album> = await res.json();
  return { listens };
}

export default function Listens() {
  const { listens: initialData } = useLoaderData<{ listens: PaginatedResponse<Listen> }>();

  return (
    <ChartLayout
      title="Last Played"
      initialData={initialData}
      endpoint="listens"
      render={({ data, page, onNext, onPrev }) => (
        <div className="flex flex-col gap-5">
        <div className="flex gap-15 mx-auto">
          <button className="default" onClick={onPrev} disabled={page <= 1}>
            Prev
          </button>
          <button className="default" onClick={onNext} disabled={!data.has_next_page}>
            Next
          </button>
        </div>
            <table>
                <tbody>
                {data.items.map((item) => (
                    <tr key={`last_listen_${item.time}`}>
                        <td className="color-fg-tertiary pr-4 text-sm" title={new Date(item.time).toString()}>{timeSince(new Date(item.time))}</td>
                        <td className="text-ellipsis overflow-hidden w-[700px]">
                            <ArtistLinks artists={item.track.artists} />{' - '}
                            <Link className="hover:text-(--color-fg-secondary)" to={`/track/${item.track.id}`}>{item.track.title}</Link>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
          <div className="flex gap-15 mx-auto">
            <button className="default" onClick={onPrev} disabled={page === 0}>
              Prev
            </button>
            <button className="default" onClick={onNext} disabled={!data.has_next_page}>
              Next
            </button>
          </div>
        </div>
      )}
    />
  );
}
