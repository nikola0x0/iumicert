"use client";

import React, { useState, useEffect } from "react";
import { useServerInsertedHTML } from "next/navigation";
import { ServerStyleSheet, StyleSheetManager } from "styled-components";
import isPropValid from "@emotion/is-prop-valid";

// Custom prop validation to filter out problematic props
const shouldForwardProp = (prop: string) => {
  if (["override"].includes(prop)) {
    return false;
  }
  return isPropValid(prop);
};

export default function StyledComponentsRegistry({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isClient, setIsClient] = useState(false);

  // Only create stylesheet once with lazy initial state
  const [styledComponentsStyleSheet] = useState(() => {
    // Only create ServerStyleSheet on the server
    if (typeof window === "undefined") {
      return new ServerStyleSheet();
    }
    return null;
  });

  useEffect(() => {
    setIsClient(true);
  }, []);

  useServerInsertedHTML(() => {
    if (!styledComponentsStyleSheet) return null;
    const styles = styledComponentsStyleSheet.getStyleElement();
    styledComponentsStyleSheet.instance.clearTag();
    return <>{styles}</>;
  });

  // On client side, just use StyleSheetManager without ServerStyleSheet
  if (isClient || typeof window !== "undefined") {
    return (
      <StyleSheetManager shouldForwardProp={shouldForwardProp}>
        {children}
      </StyleSheetManager>
    );
  }

  // On server side, use ServerStyleSheet
  if (styledComponentsStyleSheet) {
    return (
      <StyleSheetManager
        sheet={styledComponentsStyleSheet.instance}
        shouldForwardProp={shouldForwardProp}
      >
        {children}
      </StyleSheetManager>
    );
  }

  return <>{children}</>;
}
