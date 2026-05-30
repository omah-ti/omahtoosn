"use client";

import React from "react";
import katex from "katex";

interface LatexProps {
  text: string;
  className?: string;
}

export default function Latex({ text, className }: LatexProps) {
  if (!text) return null;

  // Split text by '$' to identify inline math sections
  const parts = text.split("$");

  return (
    <span className={className}>
      {parts.map((part, index) => {
        // Every odd index is a math formula (i.e. was inside $)
        if (index % 2 === 1) {
          try {
            const html = katex.renderToString(part, {
              throwOnError: false,
              displayMode: false,
            });
            return (
              <span
                key={index}
                dangerouslySetInnerHTML={{ __html: html }}
              />
            );
          } catch (e) {
            return <span key={index}>${part}$</span>;
          }
        }
        // Even indices are plain text/HTML
        return <span key={index} dangerouslySetInnerHTML={{ __html: part }} />;
      })}
    </span>
  );
}
