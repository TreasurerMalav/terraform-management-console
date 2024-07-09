import React from 'react'
import PropTypes from 'prop-types'

export default function Navbar(props) {
    return (
        <>
        <div>
            <nav className="navbar navbar-expand-lg bg-body-tertiary">
                <div className="container-fluid">
                    <a className="navbar-brand" href="/"><strong>{props.title}</strong></a>
                </div>
            </nav>
        </div>
        </>
    )
}

Navbar.propTypes = { title: PropTypes.string.isRequired }

// Navbar.defaultProps = { title: "Terraform Management Console" } -- going to be deprecated

