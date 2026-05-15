import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import '@/src/index.css'
import { TooltipProvider } from '@/components/ui/tooltip'
import HomePage from "@/src/pages/HomePage";
import Dashboard from './pages/Dashboard';
import AppLayout from './layouts/AppLayout';
import NotFound from './pages/NotFound';
import ErrorPage from './pages/ErrorPage';
import ConnectPage from './pages/ConnectPage';

createRoot(document.getElementById("root")!).render(
  <StrictMode>
      <TooltipProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route element={<AppLayout />} errorElement={<ErrorPage />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/transactions" element={<div className="text-muted-foreground">Transactions coming soon</div>} />
            <Route path="/budgets" element={<div className="text-muted-foreground">Budgets coming soon</div>} />
            <Route path="/accounts" element={<div className="text-muted-foreground">Accounts coming soon</div>} />
            <Route path="/investments" element={<div className="text-muted-foreground">Investments coming soon</div>} />
            <Route path="/connect" element={<ConnectPage />} />
          </Route>
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </StrictMode>,
);
