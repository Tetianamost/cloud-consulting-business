import React, { useState } from "react";
import styled from "styled-components";
import { theme } from "../../styles/theme";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import Icon from "../ui/icon";
import { FiLock, FiUser, FiAlertCircle } from "react-icons/fi";

const LoginContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: ${theme.space[4]};
  background-color: ${theme.colors.gray100};
`;

const LoginCard = styled.div`
  width: 100%;
  max-width: 400px;
  background-color: ${theme.colors.white};
  border-radius: ${theme.borderRadius.lg};
  box-shadow: ${theme.shadows.lg};
  padding: ${theme.space[8]};
`;

const LoginHeader = styled.div`
  text-align: center;
  margin-bottom: ${theme.space[6]};
`;

const LoginTitle = styled.h1`
  font-size: ${theme.fontSizes["2xl"]};
  font-weight: ${theme.fontWeights.bold};
  color: ${theme.colors.gray900};
  margin-bottom: ${theme.space[2]};
`;

const LoginSubtitle = styled.p`
  font-size: ${theme.fontSizes.md};
  color: ${theme.colors.gray600};
`;

const LoginForm = styled.form`
  display: flex;
  flex-direction: column;
  gap: ${theme.space[4]};
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${theme.space[2]};
`;

const Label = styled.label`
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  color: ${theme.colors.gray700};
`;

const InputWrapper = styled.div`
  position: relative;
`;

const Input = styled.input`
  width: 100%;
  padding: ${theme.space[3]} ${theme.space[3]} ${theme.space[3]}
    ${theme.space[10]};
  border: 1px solid ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  transition: ${theme.transitions.normal};

  &:focus {
    outline: none;
    border-color: ${theme.colors.primary};
    box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.2);
  }
`;

const IconWrapper = styled.div`
  position: absolute;
  top: 50%;
  left: ${theme.space[3]};
  transform: translateY(-50%);
  color: ${theme.colors.gray500};
`;

const LoginButton = styled.button`
  padding: ${theme.space[3]};
  background-color: ${theme.colors.primary};
  color: ${theme.colors.white};
  border: none;
  border-radius: ${theme.borderRadius.md};
  font-size: ${theme.fontSizes.md};
  font-weight: ${theme.fontWeights.medium};
  cursor: pointer;
  transition: ${theme.transitions.normal};
  margin-top: ${theme.space[2]};

  &:hover {
    background-color: ${theme.colors.primary};
    opacity: 0.9;
  }

  &:disabled {
    background-color: ${theme.colors.gray400};
    cursor: not-allowed;
  }
`;

const ErrorMessage = styled.div`
  display: flex;
  align-items: center;
  gap: ${theme.space[2]};
  padding: ${theme.space[3]};
  background-color: ${theme.colors.danger}10;
  border: 1px solid ${theme.colors.danger}30;
  border-radius: ${theme.borderRadius.md};
  color: ${theme.colors.danger};
  font-size: ${theme.fontSizes.sm};
  margin-bottom: ${theme.space[4]};
`;

const DemoNote = styled.div`
  margin-top: ${theme.space[4]};
  padding: ${theme.space[3]};
  background-color: ${theme.colors.primary}10;
  border: 1px solid ${theme.colors.primary}30;
  border-radius: ${theme.borderRadius.md};
  color: ${theme.colors.primary};
  font-size: ${theme.fontSizes.sm};
  text-align: center;
  line-height: 1.5;
`;

const Login: React.FC = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !password) {
      setError("Please enter both username and password");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const success = await login(username, password);

      if (success) {
        // Redirect to admin dashboard
        navigate("/admin/dashboard");
      } else {
        setError("Invalid username or password");
      }
    } catch (err: any) {
      setError(err.message || "Login failed. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <LoginContainer>
      <LoginCard>
        <LoginHeader>
          <LoginTitle>Admin Sign In</LoginTitle>
          <LoginSubtitle>
            Administrator access to the cloud consulting dashboard
          </LoginSubtitle>
        </LoginHeader>

        {error && (
          <ErrorMessage>
            <Icon icon={FiAlertCircle} size={16} />
            {error}
          </ErrorMessage>
        )}

        <LoginForm onSubmit={handleSubmit}>
          <FormGroup>
            <Label htmlFor="username">Username</Label>
            <InputWrapper>
              <IconWrapper>
                <Icon icon={FiUser} size={16} />
              </IconWrapper>
              <Input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Enter your username"
                disabled={loading}
              />
            </InputWrapper>
          </FormGroup>

          <FormGroup>
            <Label htmlFor="password">Password</Label>
            <InputWrapper>
              <IconWrapper>
                <Icon icon={FiLock} size={16} />
              </IconWrapper>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                disabled={loading}
              />
            </InputWrapper>
          </FormGroup>

          <LoginButton type="submit" disabled={loading}>
            {loading ? "Signing in..." : "Sign In"}
          </LoginButton>
        </LoginForm>

        <DemoNote>
          <strong>Demo Credentials:</strong>
          <br />
          Username: admin
          <br />
          Password: cloudadmin
        </DemoNote>
      </LoginCard>
    </LoginContainer>
  );
};

export default Login;
