'use client'
import { useEffect, useState } from 'react'
import styles from './page.module.css'
import { useGeneralCtx } from './components/ctx'
import { Header } from './components/header'
import { Board } from './components/board'
import { WinnerModal } from './components/winner-modal'
import { Chat } from './components/chat'
import { LostConnectionModal } from './components/lost-connection-modal'

export default function Home() {
  const [{ game, user, ws }, setCtx] = useGeneralCtx()
  const [serverTimeOffset, setServerTimeOffset] = useState(0);
  const [findingMatch, setFindingMatch] = useState(false)
  const [gameTimer, setGameTimer] = useState(0)
  const [fetchFindMatch, setFetchFindMatch] = useState(false)

  useEffect(() => {
    const requestStart = Date.now();
    fetch('api/time').then((v) => v.json()).then(v => {
      const now = Date.now();
      const requestDelay = now - requestStart;
      setServerTimeOffset(now - v - requestDelay)
    })
  }, [])

  useEffect(() => {
    if (!user || !ws) return
    if (ws.readyState === 1) ws.send(JSON.stringify({ type: 'game/get', data: { jwt: user.jwt } }))
  }, [user, ws])

  useEffect(() => {
    if (game) setFindingMatch(false)
  }, [game])

  useEffect(() => {
    if (!game || game.winner != 0) return
    setGameTimer(Math.floor(30 + (Date.now() - (game.lastMoveTs + serverTimeOffset)) / 1000))
    const i = setInterval(() => setGameTimer((v) => {
      if (v <= 0) {
        clearInterval(i)
        return v
      }
      return v - 1
    }), 1000)
    return () => clearInterval(i)
  }, [game, serverTimeOffset])

  const cancelMatchSearch = () => {
    if (!ws || !user) return
    ws.send(JSON.stringify({ type: "game/cancel-match-search", data: { jwt: user.jwt } }))
    setCtx((v) => ({ ...v, user: v.user?.isAnonym === true ? undefined : v.user }))
  }

  let opponent = { username: "", rating: "" }
  if (game && user) {
    const u = user.id === game.playerA.userId ? game.playerB : game.playerA
    opponent = u.isAnonym ? { username: 'Anonym', rating: '-' } : { username: u.username, rating: u.rating }
  }


  useEffect(() => {
    if (game) return setFetchFindMatch(false)
    if (!fetchFindMatch || !ws) return

    const fetchMatch = () => {
      if (user) {
        ws.send(JSON.stringify({ type: "game/find-match", data: { jwt: user.jwt, isAnonym: false } }))
      } else {
        ws.send(JSON.stringify({ type: "game/find-match", data: { jwt: "", isAnonym: true } }))
      }
    }

    if (ws.readyState === ws.OPEN) {
      fetchMatch()
    } else {
      ws.addEventListener('open', fetchMatch)
    }
    ws.addEventListener('close', () => setFetchFindMatch(true))
    setFetchFindMatch(false)
  }, [ws, fetchFindMatch, user, game])

  return (
    <div className={styles['page']}>
      <Header />
      <main className={styles['main']}>
        {!game ?
          !findingMatch ?
            <button className={styles['find-btn']} onClick={() => { setFindingMatch(true); setFetchFindMatch(true) }}>
              Find match
            </button>
            :
            <button className={styles['cancel-btn']} onClick={() => { setFindingMatch(false); cancelMatchSearch() }}>
              Cancel match search
            </button>
          :
          <>
            <div className={styles['timer']}>{`${gameTimer}${!game.winner && game.currentTurnPlayerId === user?.id ? ' Your turn!!!' : ''}`}</div>
            <div className={styles['opponent']}>{`opponent: @${opponent.username} rating: ${opponent.rating}`}</div>
            <div className={styles['container']}>
              <Board />
              <Chat />
            </div>
          </>
        }
      </main>
      <WinnerModal />
      <LostConnectionModal />
    </div >
  )
}
