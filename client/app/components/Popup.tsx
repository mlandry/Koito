import React, { type PropsWithChildren, useState } from 'react';

interface Props {
  inner: React.ReactNode
  position: string
  space: number
  extraClasses?: string
  hint?: string
}

export default function Popup({ inner, position, space, extraClasses, children }: PropsWithChildren<Props>) {
    const [isVisible, setIsVisible] = useState(false);

    let positionClasses
    let spaceCSS = {}
    if (position == "top") {
        positionClasses = `top-${space} -bottom-2 -translate-y-1/2 -translate-x-1/2`
    } else if (position == "right") {
        positionClasses = `bottom-1 -translate-x-1/2`
        spaceCSS = {left: 70 + space}
    }

    return (
        <div
        className="relative"
        onMouseEnter={() => setIsVisible(true)}
        onMouseLeave={() => setIsVisible(false)}
        >
        {children}
        <div
            className={`
            absolute 
            ${positionClasses}
            ${extraClasses ? extraClasses : ''}
            bg-(--color-bg) color-fg border-1 border-(--color-bg-tertiary)
            px-3 py-2 rounded-lg
            transition-opacity duration-100
            ${isVisible ? 'opacity-100' : 'opacity-0 pointer-events-none'}
            z-50 text-center
            flex
            `}
            style={spaceCSS}
        >
            {inner}
        </div>
        </div>
    );
}
