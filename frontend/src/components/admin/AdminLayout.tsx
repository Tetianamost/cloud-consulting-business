import React from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { FiHome, FiList, FiBarChart2, FiMail, FiLogOut } from 'react-icons/fi';
import { theme } from '../../styles/theme';
import Icon from '../ui/Icon';

const AdminContainer = styled.div`
  display: flex;
  min-height: 100vh;
`;

const Sidebar = styled.aside`
  width: 250px;
  background-color: ${theme.colors.gray900};
  color: ${theme.colors.white};
  padding: ${theme.space[4]};
  display: flex;
  flex-direction: column;
  
  @media (max-width: ${theme.breakpoints.md}) {
    width: 80px;
    padding: ${theme.space[3]};
  }
`;

const Logo = styled.div`
  font-family: ${theme.fonts.heading};
  font-weight: ${theme.fontWeights.bold};
  font-size: ${theme.fontSizes.xl};
  color: ${theme.colors.white};
  margin-bottom: ${theme.space[6]};
  padding: ${theme.space[2]} 0;
  border-bottom: 1px solid ${theme.colors.gray700};
  
  span {
    color: ${theme.colors.secondary};
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    font-size: ${theme.fontSizes.md};
    text-align: center;
  }
`;

const NavMenu = styled.nav`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

const NavItem = styled(NavLink)`
  display: flex;
  align-items: center;
  padding: ${theme.space[3]};
  color: ${theme.colors.gray300};
  border-radius: ${theme.borderRadius.md};
  margin-bottom: ${theme.space[2]};
  text-decoration: none;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.gray800};
    color: ${theme.colors.white};
  }
  
  &.active {
    background-color: ${theme.colors.primary};
    color: ${theme.colors.white};
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    justify-content: center;
    padding: ${theme.space[3]} ${theme.space[2]};
  }
`;

const NavText = styled.span`
  margin-left: ${theme.space[3]};
  
  @media (max-width: ${theme.breakpoints.md}) {
    display: none;
  }
`;

const LogoutButton = styled.button`
  display: flex;
  align-items: center;
  padding: ${theme.space[3]};
  color: ${theme.colors.gray300};
  background: none;
  border: none;
  border-radius: ${theme.borderRadius.md};
  margin-top: auto;
  cursor: pointer;
  transition: ${theme.transitions.normal};
  
  &:hover {
    background-color: ${theme.colors.gray800};
    color: ${theme.colors.white};
  }
  
  @media (max-width: ${theme.breakpoints.md}) {
    justify-content: center;
    padding: ${theme.space[3]} ${theme.space[2]};
  }
`;

const Content = styled.main`
  flex: 1;
  padding: ${theme.space[6]};
  background-color: ${theme.colors.light};
  overflow-y: auto;
`;

const AdminLayout: React.FC = () => {
  const navigate = useNavigate();
  
  const handleLogout = () => {
    // For demo purposes, just navigate to home
    window.location.href = '/';
  };
  
  return (
    <AdminContainer>
      <Sidebar>
        <Logo>
          Cloud<span>Admin</span>
        </Logo>
        <NavMenu>
          <NavItem to="/admin" end>
            <Icon icon={FiHome} size={20} />
            <NavText>Dashboard</NavText>
          </NavItem>
          <NavItem to="/admin/inquiries">
            <Icon icon={FiList} size={20} />
            <NavText>Inquiries</NavText>
          </NavItem>
          <NavItem to="/admin/metrics">
            <Icon icon={FiBarChart2} size={20} />
            <NavText>Metrics</NavText>
          </NavItem>
          <NavItem to="/admin/email-status">
            <Icon icon={FiMail} size={20} />
            <NavText>Email Status</NavText>
          </NavItem>
        </NavMenu>
        <LogoutButton onClick={handleLogout}>
          <Icon icon={FiLogOut} size={20} />
          <NavText>Logout</NavText>
        </LogoutButton>
      </Sidebar>
      <Content>
        <Outlet />
      </Content>
    </AdminContainer>
  );
};

export default AdminLayout;