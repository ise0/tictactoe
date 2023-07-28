import styles from "./winner.module.css"
import { Modal } from "../modal"
import { useGeneralCtx } from "../ctx"

const winnerMods: Record<string, string> = {
    yes: styles['winner_yes'],
    no: styles['winner_no']
}

export function WinnerModal() {
    const [{ game, user }, setCtx] = useGeneralCtx()
    if (!game || !user || !game.winner || game.winner === 0) {
        return <></>
    }
    const win = game.winner === user.id;

    const close = () => setCtx(v => ({
        ...v,
        game: undefined,
        user: v.user?.isAnonym ? undefined : v.user
    }
    ))

    return (
        <Modal alignX="center" alignY="center" autoClose={close}>
            <div className={`${styles["winner"]} ${winnerMods[win ? "yes" : "no"]}`} onClick={close}>
                <span className={styles["text"]}>{win ? "you win! :)" : "you lose :("}</span>
            </div>
        </Modal >
    )
}