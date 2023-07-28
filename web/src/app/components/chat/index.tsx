import { useState } from "react";
import { useGeneralCtx } from "../ctx";
import styles from "./chat.module.css"

export function Chat() {
    const [{ game, user, ws }] = useGeneralCtx()
    const [msg, setMsg] = useState("")

    function addChatMsg() {
        if (!ws || !user) return
        ws.send(JSON.stringify({
            type: "game/add-chat-msg",
            data: {
                jwt: user.jwt,
                text: msg
            }
        }))
        setMsg("")
    }

    return (
        <div className={styles['chat']}>
            <ol className={styles['list']}>
                {game?.chat.map(el =>
                    <li className={`${styles['item']} ${styles[user?.id === el.userId ? 'item_me' : 'item_opponent']}`} key={el.ts}>
                        <div className={styles['msg']}>{el.text}</div>
                        <div className={styles['ts']}>{new Date(el.ts).toLocaleTimeString().slice(0, 5)}</div>
                    </li>
                )}
            </ol>
            <div className={styles['msg-input']}>
                <input
                    className={styles['msg-ctrl']}
                    type="text"
                    value={msg}
                    onChange={(e) => setMsg(e.target.value)}
                    onKeyDown={(evt) => {
                        if (evt.key === 'Enter') {
                            addChatMsg()
                        }
                    }}
                />
                <button className={styles['btn-send']} onClick={addChatMsg}>send</button>
            </div>
        </div>
    )
}