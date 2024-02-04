import React from "react";
import Container from "@mui/material/Container";

interface RegistrationCompletedScreenProps {
    email: string
}

const RegistrationCompletedScreen: React.FC<RegistrationCompletedScreenProps> = (props) => {
    return (
        <Container style={{ width: "600px", height: "400px", display: "flex", flexDirection: "column", justifyContent: "space-evenly" }}>
            <h1 style={{color:"#11BB00"}}>Congratulations!</h1>
            <h4>Email address {props.email} has been successfully authorized!</h4>
        </Container>
    )
}

export default RegistrationCompletedScreen