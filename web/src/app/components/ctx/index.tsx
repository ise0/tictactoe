'use client'

import { Dispatch, PropsWithChildren, SetStateAction, createContext, useContext, useEffect, useState } from "react";
import { wsReducer } from "./ws-reducer";

type Msg = {
    userId: number,
    ts: number,
    text: string
}

type Player = {
    userId: number,
    isAnonym: boolean,
    username: string,
    playerSign: string,
    rating: string,
}

type Game = {
    gameSessionId: number,
    board: string[][],
    winner: number,
    chat: Msg[],
    currentTurnPlayerId: number,
    lastMoveTs: number,
    playerA: Player,
    playerB: Player
}

type User = { id: number; username: string, isAnonym: boolean, rating: number, jwt: string }

export type Ctx = { user?: User, game?: Game, ws?: WebSocket }

const ctx = createContext<[Ctx, Dispatch<SetStateAction<Ctx>>] | undefined>(undefined)

export function GeneralCtx({ children }: PropsWithChildren) {
    const v = useState<Ctx>({});
    const [value, setValue] = v;

    useEffect(() => {
        function reducer(ev: MessageEvent<any>) {
            setValue((v) => wsReducer(v, ev))
        }

        function connect() {
            const newWS = new WebSocket("ws://" + window.location.host + "/api/ws");
            setValue((v) => ({ ...v, ws: newWS }))
            newWS.addEventListener("message", reducer)
            newWS.addEventListener("close", (e) => {
                setValue((v) => ({ ...v, ws: undefined }))
                setTimeout(connect, 1000)
            })
        }
        connect()
    }, [setValue])

    return (
        <ctx.Provider value={v}>
            {children}
        </ctx.Provider>
    )
}

export function useGeneralCtx() {
    const v = useContext(ctx)
    if (!v) {
        throw new Error("useUser must be used within a UserCtx");
    }
    return v
}