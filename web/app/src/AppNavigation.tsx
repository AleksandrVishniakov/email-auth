import React, {useState} from "react";
import EmailVerificationScreen from "./components/email-verification/EmailVerificationScreen";
import EmailInputScreen from "./components/email-input/EmailInputScreen";
import RegistrationCompletedScreen from "./components/registration-completed/RegistrationCompletedScreen";

enum Screens {
    Login,
    Verification,
    RegistrationCompleted
}

interface AppNavigationProps {
    onEmailSubmit: (email: string) => Promise<boolean>
    onVerificationSubmit: (code: string) => Promise<void>
    errorHandler: (text: string) => void
}
const AppNavigation: React.FC<AppNavigationProps> = (props) => {
    const [currentScreen, setCurrentScreen] = useState(Screens.Login)
    const [userEmail, setUserEmail] = useState("")
    const [userAuthorized, setUserAuthorized] = useState(false)

    switch (currentScreen) {
        case Screens.Verification:
            return <EmailVerificationScreen
                isUserAuthorized={userAuthorized}
                email={userEmail}
                onCodeSubmit={ async (code: string) => {
                    let ok = true

                    try {
                        await props.onVerificationSubmit(code)
                    }
                    catch (e: any) {
                        ok = false
                        props.errorHandler(e.toString())
                    }

                    if (ok) {
                        setCurrentScreen(Screens.RegistrationCompleted)
                    }
                }}
            />
        case Screens.RegistrationCompleted:
            return <RegistrationCompletedScreen email={userEmail}/>

        default:
            return <EmailInputScreen onSubmit={
                async (email: string) => {
                    setUserEmail(email)
                    let ok = true

                    try {
                        const isUserAuthorized = await props.onEmailSubmit(email)
                        setUserAuthorized(isUserAuthorized)
                    }
                    catch (e: any) {
                        ok = false
                        props.errorHandler(e.toString())
                    }

                    if (ok) {
                        setCurrentScreen(Screens.Verification)
                    }
                }
            }/>
    }
}

export default AppNavigation