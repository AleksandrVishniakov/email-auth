import Container from "@mui/material/Container";
import React, {useState} from "react";
import {FormControl, Input, InputAdornment, InputLabel, Link} from "@mui/material";
import Button from "@mui/material/Button";

interface EmailVerificationScreenProps {
    email: string
    isUserAuthorized: boolean
    onCodeSubmit: (code: string) => void
}

const EmailVerificationScreen: React.FC<EmailVerificationScreenProps> = (props) => {
    const [codeInputValue, setCodeInputValue] = useState("")

    const handleCodeInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setCodeInputValue(event.target.value)
    }
    const handleCodeSubmit = () => {
        props.onCodeSubmit(codeInputValue)
    }

    return (
        <Container style={{
            width: "600px",
            height: "400px",
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-evenly"
        }}>
            <div>
                <h1>Email Verification</h1>
                <h3>{props.isUserAuthorized ? "Sign In" : "Sign Up"}</h3>
                <p>{props.isUserAuthorized ? "Glad to see you again!" : "Nice to meet you, new user!"}</p>
                <p>We send an email to <Link href="#">{props.email}</Link>. Please, enter the provided 6-digit code
                    below</p>
            </div>

            <FormControl sx={{m: 1}} variant="standard"
                         style={{
                             width: "80%",
                             marginLeft: "auto",
                             marginRight: "auto",
                         }}
            >
                <InputLabel htmlFor="standard-adornment-amount">Verification Code</InputLabel>
                <Input
                    id="standard-adornment-code"
                    startAdornment={<InputAdornment position="start">#</InputAdornment>}
                    style={{
                        fontSize: "32px",
                        letterSpacing: "50px"
                    }}
                    inputProps={{
                        maxLength: 6,
                        inputMode: 'numeric',
                        pattern: '[0-9]*'
                    }}
                    type="text"

                    onChange={handleCodeInputChange}
                />
            </FormControl>

            <Button variant="contained" onClick={handleCodeSubmit}>Verify</Button>
        </Container>
    )
}

export default EmailVerificationScreen