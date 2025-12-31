'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { ComponentProps, useCallback, useTransition } from 'react';

type TransitionLinkProps = ComponentProps<typeof Link> & {
  onClick?: () => void;
};

export function TransitionLink({ href, onClick, children, ...props }: TransitionLinkProps) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const handleClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      onClick?.();
      
      const url = typeof href === 'string' ? href : href.pathname || '/';
      
      if (document.startViewTransition) {
        document.startViewTransition(() => {
          startTransition(() => {
            router.push(url);
          });
        });
      } else {
        startTransition(() => {
          router.push(url);
        });
      }
    },
    [href, onClick, router]
  );

  return (
    <Link href={href} onClick={handleClick} {...props}>
      {children}
    </Link>
  );
}
