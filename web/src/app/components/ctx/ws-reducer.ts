import type { Ctx } from "./index"

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
    currentTurnPlayerId: number,
    lastMoveTs: number,
    chat: Msg[],
    playerA: Player,
    playerB: Player
}


type Actions = {
    type: "game/game",
    data: Game
} | {
    type: "game/gameBoard",
    data: { board: string[][], currentTurnPlayerId: number, winner: number }
} | {
    type: "game/timeExpired",
    data: { winner: number }
} | {
    type: "game/chatMsg",
    data: Msg
} | {
    type: "game/anonJwt",
    data: { jwt: string, userId: number }
}


export function wsReducer(data: Ctx, ev: MessageEvent<any>): Ctx {
    const a: Actions = JSON.parse(ev.data)
    switch (a.type) {
        case 'game/game':
            return { ...data, game: { ...a.data, winner: 0 } }
        case 'game/gameBoard': {
            let user = data.user
            if (a.data.winner != 0 && user && !user.isAnonym) {
                let rating = user.rating + (user.id === a.data.winner ? 100 : -100);
                rating = rating < 0 ? 0 : rating
                user = { ...user, rating }
            }
            return { ...data, game: data.game ? { ...data.game, ...a.data } : undefined, user }
        }
        case 'game/timeExpired': {
            let user = data.user
            if (a.data.winner != 0 && user && !user.isAnonym) {
                let rating = user.rating + (user.id === a.data.winner ? 100 : -100);
                rating = rating < 0 ? 0 : rating
                user = { ...user, rating }
            }
            return { ...data, game: data.game ? { ...data.game, ...a.data } : undefined, user }
        }
        case 'game/chatMsg':
            if (!data.game) return data
            const newChat = [...data.game.chat, { ...a.data, ts: Date.now() }]
            return ({ ...data, game: { ...data.game, chat: newChat } })
        case 'game/anonJwt':
            return ({ ...data, user: { isAnonym: true, id: a.data.userId, jwt: a.data.jwt, rating: 0, username: "Anonym" } })
    }

    return data
}