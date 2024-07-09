import React from 'react'
import PropTypes from 'prop-types'
import Navbar from './Navbar'
import TerraformConsole from './TerraformConsole';

export default function Home() {

    return (
        <>
        <div className="container">
            <Navbar title="Terraform Management Console"/>
        </div>
        <div className="container">
            <TerraformConsole />
        </div>
        </>
    )
}

Navbar.propTypes = { title: PropTypes.string.isRequired }

// Navbar.defaultProps = { title: "Terraform Management Console" } -- going to be deprecated

