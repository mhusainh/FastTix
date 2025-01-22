// src/app/components/Footer.tsx

import Link from 'next/link';

export default function Footer() {
  return (
    <footer className="bg-gray-800 text-gray-300 p-4 text-center">
      <div className="mb-2">
        © 2024 TiketKu. All rights reserved.
      </div>
      <div className="flex justify-center gap-4">
        <Link href="/privacy" className="hover:underline">Privacy Policy</Link>
        <Link href="/terms" className="hover:underline">Terms of Service</Link>
        <Link href="/contact" className="hover:underline">Contact</Link>
      </div>
    </footer>
  );
}
