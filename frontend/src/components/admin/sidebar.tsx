import React from "react"
import { Link, useLocation } from "react-router-dom"
import { BarChart3, Mail, MessageSquare, Settings, FileText, Bot } from "lucide-react"
import styled from "styled-components"
import { Button } from "../ui/button"
import { theme } from "../../styles/theme"

const SidebarContainer = styled.div`
  display: none;
  border-right: 1px solid ${theme.colors.gray200};
  background-color: ${theme.colors.gray100};
  width: 256px;
  
  @media (min-width: ${theme.breakpoints.lg}) {
    display: block;
  }
`

const SidebarContent = styled.div`
  display: flex;
  height: 100vh;
  max-height: 100vh;
  flex-direction: column;
  gap: 8px;
`

const Header = styled.div`
  display: flex;
  height: 60px;
  align-items: center;
  border-bottom: 1px solid ${theme.colors.gray200};
  padding: 0 16px;
`

const HeaderLink = styled(Link)`
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: ${theme.fontWeights.semibold};
  color: ${theme.colors.primary};
  text-decoration: none;
  
  &:hover {
    color: ${theme.colors.secondary};
  }
`

const NavContainer = styled.div`
  flex: 1;
  overflow: auto;
  padding: 8px 0;
`

const Nav = styled.nav`
  display: grid;
  align-items: start;
  padding: 0 8px;
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
`

const NavButton = styled(Button)`
  justify-content: flex-start;
  margin-bottom: 4px;
`

const Footer = styled.div`
  margin-top: auto;
  padding: 16px;
`

const UserInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: ${theme.borderRadius.lg};
  background-color: ${theme.colors.gray200};
  padding: 16px;
`

const UserDetails = styled.div`
  display: grid;
  gap: 4px;
`

const UserName = styled.p`
  font-size: ${theme.fontSizes.sm};
  font-weight: ${theme.fontWeights.medium};
  line-height: 1;
  margin: 0;
`

const UserEmail = styled.p`
  font-size: ${theme.fontSizes.xs};
  color: ${theme.colors.gray600};
  margin: 0;
`

const navItems = [
  {
    title: "Dashboard",
    href: "/admin/dashboard",
    icon: BarChart3,
  },
  {
    title: "Inquiries",
    href: "/admin/inquiries",
    icon: MessageSquare,
  },
  {
    title: "AI Chat",
    href: "/admin/chat",
    icon: Bot,
  },
  {
    title: "Simple Chat",
    href: "/admin/simple-chat",
    icon: MessageSquare,
  },
  {
    title: "Chat Mode",
    href: "/admin/chat-mode",
    icon: Settings,
  },
  {
    title: "Metrics",
    href: "/admin/metrics",
    icon: FileText,
  },
  {
    title: "Email Status",
    href: "/admin/email-status",
    icon: Mail,
  },
]

export function AdminSidebar() {
  const location = useLocation()
  const pathname = location.pathname

  return (
    <SidebarContainer>
      <SidebarContent>
        <Header>
          <HeaderLink to="/admin">
            <Bot size={24} />
            <span>AI Admin Portal</span>
          </HeaderLink>
        </Header>
        <NavContainer>
          <Nav>
            {navItems.map((item) => (
              <NavButton
                key={item.href}
                asChild
                variant={pathname === item.href ? "secondary" : "ghost"}
              >
                <Link to={item.href}>
                  <item.icon size={16} style={{ marginRight: '8px' }} />
                  {item.title}
                </Link>
              </NavButton>
            ))}
          </Nav>
        </NavContainer>
        <Footer>
          <UserInfo>
            <UserDetails>
              <UserName>Admin User</UserName>
              <UserEmail>admin@example.com</UserEmail>
            </UserDetails>
          </UserInfo>
        </Footer>
      </SidebarContent>
    </SidebarContainer>
  )
}
