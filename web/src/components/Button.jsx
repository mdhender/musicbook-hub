// src/components/Button.jsx
// Copyright (c) 2025 Michael D Henderson. All rights reserved.

import React from "react";
import classNames from "classnames";

const baseClasses = "px-4 py-2 rounded font-medium focus:outline-none focus:ring-2 focus:ring-offset-2";

const variants = {
    blue: "bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-400",
    red: "bg-red-600 text-white hover:bg-red-700 focus:ring-red-400",
    ghost: "bg-transparent text-gray-700 hover:underline focus:ring-gray-400",
    black: "bg-black text-white hover:bg-gray-800 focus:ring-gray-600",
    disabled: "bg-gray-300 text-gray-500 cursor-not-allowed",
};

export default function Button({
                                   variant = "blue",
                                   className = "",
                                   disabled = false,
                                   ...props
                               }) {
    const classes = classNames(
        baseClasses,
        variants[disabled ? "disabled" : variant],
        className
    );

    return <button className={classes} disabled={disabled} {...props} />;
}
