import { useState } from "react"
import { useGeneralCtx } from "../ctx"
import styles from "./board.module.css"

const cellSignMods: Record<string, string> = {
    X: styles['cell_sign_x'],
    O: styles['cell_sign_o']
}
const cellHoverMods: Record<string, string> = {
    X: styles['cell_hover_x'],
    O: styles['cell_hover_o'],
}

export function Board() {
    const [{ game, user, ws }] = useGeneralCtx()
    const [hoverCell, setHoverCell] = useState<{ x: number, y: number }>()
    if (!user || !game) return <></>
    function updateMatch(coords: { x: number, y: number }) {
        if (!ws || !user) return
        ws?.send(JSON.stringify({ type: "game/do-move", value: { jwt: user.jwt, coords } }))
    }

    const myTurn = game.currentTurnPlayerId === user.id;
    const mySign = user.id === game.playerA.userId ? game.playerA.playerSign : game.playerB.playerSign;

    return (
        <div className={styles['container-2']}>
            <table className={styles['board']}>
                <tbody>
                    {game?.board.map((el, i) =>
                        <tr key={i} className={styles['row']}>
                            {el.map((el2, i2) =>
                                <td key={i2}
                                    onPointerEnter={() => setHoverCell({ x: i2, y: i })}
                                    onPointerLeave={() => setHoverCell(undefined)}
                                    className={`${styles['cell']} ${cellSignMods[el2] || (myTurn && hoverCell?.x === i2 && hoverCell.y === i ? cellHoverMods[mySign] : '')}`}
                                    onClick={() => updateMatch({ x: i2, y: i })}
                                />
                            )}
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}