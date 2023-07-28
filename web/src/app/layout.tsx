import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { GeneralCtx } from './components/ctx'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'TicTacToe'
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <GeneralCtx>
      <html lang="en">
        <body className={inter.className}>{children}</body>
      </html>
    </GeneralCtx>
  )
}
