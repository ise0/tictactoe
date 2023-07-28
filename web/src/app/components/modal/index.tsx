'use client'
import { ReactNode, useEffect, useRef, useState } from 'react';
import ReactDOM from 'react-dom';
import styles from './modal.module.css';

type Props = {
  children: ReactNode;
  autoClose?: () => void;
  alignX?: 'left' | 'center' | 'right';
  alignY?: 'top' | 'center' | 'bottom';
};

const AlignY = {
  top: 'top',
  center: 'middle',
  bottom: 'bottom',
};

export function Modal({ children, autoClose, alignX, alignY }: Props) {
  const modalRef = useRef<HTMLDivElement>(null);
  const [bodyElem, setBodyElem] = useState<HTMLElement | null>(null);
  useEffect(() => {
    const newBodyElem = document.body;
    if (!newBodyElem) throw new Error();
    setBodyElem(newBodyElem);
  }, []);

  useEffect(() => {
    const onOutsideClick = (evt: Event) => {
      if (autoClose && modalRef.current && modalRef.current === evt.target) autoClose();
    };
    const onEsc = (evt: KeyboardEvent) => {
      if (autoClose && evt.key === 'Escape') autoClose();
    };
    document.addEventListener('click', onOutsideClick, true);
    document.addEventListener('keydown', onEsc, true);
    return () => {
      document.removeEventListener('click', onOutsideClick, true);
      document.removeEventListener('keydown', onEsc, true);
    };
  }, [autoClose]);

  useEffect(() => {
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = '';
    };
  }, []);

  useEffect(() => {
    const outerFocusElementRef = document.activeElement;
    modalRef.current?.focus({ preventScroll: true });
    return () => {
      if (outerFocusElementRef) (outerFocusElementRef as HTMLElement).focus();
    };
  }, []);

  if (!bodyElem) return null;

  return ReactDOM.createPortal(
    <div className={styles['modal']} ref={modalRef} style={{ textAlign: alignX }} tabIndex={0}>
      <div className={styles['modal-inner']} style={{ verticalAlign: alignY && AlignY[alignY] }}>{children}</div>
      {autoClose && (
        <div
          className="visually-hidden"
          tabIndex={0}
          onFocus={autoClose}
          onBlur={(evt) => evt.preventDefault()}
        />
      )}
    </div>,
    bodyElem
  );
}
