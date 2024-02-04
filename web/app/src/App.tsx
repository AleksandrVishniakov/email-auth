import './App.css';
import AppNavigation from "./AppNavigation";
import Snackbar from '@mui/material/Snackbar';
import MuiAlert, {AlertProps} from '@mui/material/Alert';
import React, {forwardRef, useState} from "react";

interface APIError {
    message: string
    code: number
    timestamp: Date
}

function App() {
    const [errorText, setErrorText] = useState("")
    const [snackbarOpen, setSnackbarOpen] = useState(false)
    const [userEmail, setUserEmail] = useState("")

    const handleSnackbarClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setSnackbarOpen(false);
    };

    const openErrorSnackbar = (text: string) => {
        setErrorText(text)
        setSnackbarOpen(true)
    }

    return (
        <div className="App">
            <AppNavigation
                errorHandler={openErrorSnackbar}

                onEmailSubmit={
                    async (email: string): Promise<boolean> => {
                        setUserEmail(email)

                        if (email.length === 0) {
                            openErrorSnackbar("field cannot be empty")
                            throw new Error("field cannot be empty")
                        }

                        return authRequest(email)
                    }

                }
                onVerificationSubmit={
                    async (code: string) => {
                        if (code.length !== 6) {
                            openErrorSnackbar("invalid code length")
                            throw new Error("invalid code length")
                        }

                        return emailVerificationRequest(userEmail, code)
                    }
                }
            />

            <Snackbar open={snackbarOpen} autoHideDuration={2000} onClose={handleSnackbarClose}
                      anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}>
                <Alert onClose={handleSnackbarClose} severity="error" sx={{width: '100%'}}>
                    {errorText}
                </Alert>
            </Snackbar>
        </div>
    );
}

const authRequest = async (email: string) => {
    interface authResponse {
        isUserAuthorized: boolean
    }

    const response = await fetch("http://localhost:8003/auth?email=" + email)

    if (!response.ok) {
        const json = await response.json()
        const apiError = json as APIError
        throw new Error(apiError.code + ". " + apiError.message)
    }

    const json = await response.json()
    const apiResponse = json as authResponse

    return apiResponse.isUserAuthorized
}

const emailVerificationRequest = async (email: string, code: string) => {
    const response = await fetch("http://localhost:8003/verify/" + email + "?code=" + code)

    if (!response.ok) {
        const json = await response.json()
        const apiError = json as APIError
        throw new Error(apiError.code + ". " + apiError.message)
    }
}


const Alert = forwardRef<HTMLDivElement, AlertProps>(function Alert(
    props,
    ref,
) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default App;
