import React, { useState } from "react"
import Container from '@mui/material/Container';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';

interface EmailInputScreenProps {
    onSubmit: (email: string) => void
}

const EmailInputScreen: React.FC<EmailInputScreenProps> = (props) => {
    const [emailInputValue, setEmailInputValue] = useState("")

    const handleInputChange = (evt : React.ChangeEvent<HTMLInputElement>) => {
        setEmailInputValue(evt.target.value)
    }

    const handleEmailSubmit = () => {
        props.onSubmit(emailInputValue)
    }

    return (
        <Container style={{ width: "600px", height: "400px", display: "flex", flexDirection: "column", justifyContent: "space-evenly" }}>

            <div>
                <h1>Login</h1>
                <p>Enter your email to create new account or login to an existing one</p>
            </div>

            <TextField
                required
                id="outlined-required"
                label="Email"
                type="email"
                fullWidth
                value={emailInputValue}
                onChange={handleInputChange}
            />

            <Button variant="contained" onClick={handleEmailSubmit}>Login</Button>

        </Container>
    )
}

export default EmailInputScreen