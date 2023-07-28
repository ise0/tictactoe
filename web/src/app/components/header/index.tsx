import { useState } from 'react'
import { useGeneralCtx } from '../ctx'
import styles from './header.module.css'
import { AuthModal } from '../auth-modal'


export function Header() {
  const [showSignUp, setShowSignUp] = useState(false)
  const [showSignIn, setShowSignIn] = useState(false)
  const [{ user, ws }, setCtx] = useGeneralCtx()

  return (
    <header className={styles['header']}>
      {user && !user.isAnonym && <span className={styles["user-rating"]}>rating: {user.rating}</span>}
      <div className={styles['container']}>
        {user && !user.isAnonym ?
          <>
            <span className={styles["username"]}>@{user.username}</span>
            <button
              className={styles["button"]}
              onClick={() => { 
                setCtx({ user: undefined, game: undefined });
                ws?.close() 
              }}
            >
              Logout
            </button>
          </> :
          <>
            <button className={styles["button"]} onClick={() => setShowSignIn(true)}>Sign in</button> /
            <button className={styles["button"]} onClick={() => setShowSignUp(true)}>Sign up</button>
          </>
        }
      </div>
      {showSignUp && <AuthModal closeModal={() => setShowSignUp(false)} action='Sign Up' />}
      {showSignIn && <AuthModal closeModal={() => setShowSignIn(false)} action='Sign In' />}
    </header >
  )
}