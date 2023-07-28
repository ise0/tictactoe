'use client'

import { useState } from "react";
import { Modal } from "../modal";
import styles from "./auth.module.css"
import { useGeneralCtx } from "../ctx";


type Props = {
    closeModal: () => void,
    action: "Sign In" | "Sign Up"
};

export function AuthModal({ closeModal, action }: Props) {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [, setCtx] = useGeneralCtx()

    const onSubmit = () => {
        fetch(`api/user/${action === 'Sign In' ? 'sign-in' : 'sign-up'}`,
            {
                method: "POST",
                body: JSON.stringify({ username, password }),
                headers: { 'Content-Type': 'application/json' }
            }
        ).then((v) => {
            if (v.ok) {
                return v.json()
            }
        }).then(v => {
            if (v != undefined) {
                setCtx((prevValue) => ({ ...prevValue, user: v, game: undefined }))
                closeModal()
            }
        })
    }

    return (
        <Modal alignX="center" alignY="center" autoClose={closeModal}>
            <div className={styles["auth"]}>
                <h3 className={styles["title"]}>{action}</h3>
                <label className={styles['input']}>
                    <span className={styles['input__label']}>username</span>
                    <input
                        className={styles['input__ctrl']}
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        onKeyDown={(evt) => {
                            if (evt.key === 'Enter') {
                                onSubmit()
                            }
                        }}
                    />
                </label>
                <label className={styles['input']}>
                    <span className={styles['input__label']}>password</span>
                    <input
                        className={styles['input__ctrl']}
                        type="text"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        onKeyDown={(evt) => {
                            if (evt.key === 'Enter') {
                                onSubmit()
                            }
                        }}
                    />
                </label>
                <button onClick={onSubmit} className={styles['submit']}>Submit</button>
            </div>
        </Modal>
    )
}