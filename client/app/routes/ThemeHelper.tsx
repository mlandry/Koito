import { useState } from "react"
import { useAppContext } from "~/providers/AppProvider"
import { AsyncButton } from "../components/AsyncButton"
import AllTimeStats from "~/components/AllTimeStats"
import ActivityGrid from "~/components/ActivityGrid"
import LastPlays from "~/components/LastPlays"
import TopAlbums from "~/components/TopAlbums"
import TopArtists from "~/components/TopArtists"
import TopTracks from "~/components/TopTracks"

export default function ThemeHelper() {

    const homeItems = 3

    return (
        <div className="mt-10 flex flex-col gap-10 items-center">
            <div className="flex gap-5">
                <AllTimeStats />
                <ActivityGrid />
            </div>
            <div className="flex flex-wrap 2xl:gap-20 xl:gap-10 justify-around gap-5">
                <TopArtists period="all_time" limit={homeItems} />
                <TopAlbums period="all_time" limit={homeItems} />
                <TopTracks period="all_time" limit={homeItems} />
                <LastPlays limit={Math.floor(homeItems * 2.5)} />
            </div>
            <div className="flex flex-col gap-6 bg-secondary p-10 rounded-lg">
                <div className="flex flex-col gap-4 items-center">
                    <p>You're logged in as <strong>Example User</strong></p>
                    <AsyncButton loading={false} onClick={() => {}}>Logout</AsyncButton>
                </div>
                <div className="flex flex gap-4">
                    <input
                        name="koito-update-username"
                        type="text"
                        placeholder="Update username"
                        className="w-full mx-auto fg bg rounded p-2"
                    />
                    <AsyncButton loading={false} onClick={() => {}}>Submit</AsyncButton>
                </div>
                <div className="flex flex gap-4">
                    <input
                        name="koito-update-password"
                        type="password"
                        placeholder="Update password"
                        className="w-full mx-auto fg bg rounded p-2"
                    />
                    <input
                        name="koito-confirm-password"
                        type="password"
                        placeholder="Confirm password"
                        className="w-full mx-auto fg bg rounded p-2"
                    />
                    <AsyncButton loading={false} onClick={() => {}}>Submit</AsyncButton>
                </div>
                <div className="flex gap-2 mt-3">
                    <input type="checkbox" name="reverse-merge-order" onChange={() => {}} />
                    <label htmlFor="reverse-merge-order">Example checkbox</label>
                </div>
                <p className="success">successfully displayed example text</p>
                <p className="error">this is an example of error text</p>
                <p className="info">here is an informational example</p>
                <p className="warning">heed this warning, traveller</p>
            </div>
        </div>
    )
}