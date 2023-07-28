import styles from "./lost-conn.module.css"
import { Modal } from "../modal"
import { useGeneralCtx } from "../ctx"
import { useEffect, useState } from "react";

export function LostConnectionModal() {
    const [{ ws }] = useGeneralCtx();
    const [show, setShow] = useState(true);

    useEffect(() => {
        if (!ws || ws.readyState === ws.CLOSED) {
            setShow(true)
            if (ws) {
                const setShowTrue = () => setShow(false)
                ws.addEventListener('open', setShowTrue)
                return () => ws.removeEventListener('open', setShowTrue)
            }
        } else {
            setShow(false)
        }
    }, [ws])

    if (!show) return <></>
    return (
        <Modal alignX="center" alignY="center">
            <div className={styles['lost-conn']}>lost connection to the server. Trying reconnect...</div>
        </Modal >
    )
}