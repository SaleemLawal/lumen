import {Button} from "@/components/ui/button"
import { Link } from "react-router-dom"
import {
  BarChart3,
  PiggyBank,
  ShieldCheck,
  TrendingUp,
  Wallet,
  ArrowRight,
  Sparkles,
} from "lucide-react";

const features = [
  {
    icon: BarChart3,
    title: "Net worth, at a glance",
    desc: "A clean timeline of how your wealth is moving across any range you care about.",
  },
  {
    icon: Wallet,
    title: "All accounts, one view",
    desc: "Cash, credit, investments and loans grouped just the way you think about them.",
  },
  {
    icon: PiggyBank,
    title: "Budgets that don't nag",
    desc: "Set thoughtful limits and see your progress without the guilt-trip.",
  },
  {
    icon: TrendingUp,
    title: "Investment tracking",
    desc: "Holdings, allocations and performance; Without leaving your dashboard.",
  },
  {
    icon: ShieldCheck,
    title: "Private by default",
    desc: "Your data is yours. No upsells, no ad targeting, ever.",
  },
  {
    icon: Sparkles,
    title: "Built for focus",
    desc: "A neutral, distraction-free interface that respects your attention.",
  },
];

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="sticky top-0 z-30 border-b border-border/60 bg-background/80 backdrop-blur-xl">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 md:px-6">
          <Link to="/" className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-foreground text-background">
              <Sparkles className="h-4 w-4" />
            </div>
            <span className="text-base font-semibold tracking-tight">
              Lumen
            </span>
          </Link>
          <nav className="hidden items-center gap-7 text-sm text-muted-foreground md:flex">
            <a href="#features" className="hover:text-foreground">
              Features
            </a>
            <a href="#pricing" className="hover:text-foreground">
              Pricing
            </a>
            <a href="#faq" className="hover:text-foreground">
              FAQ
            </a>
          </nav>
          <div className="flex items-center gap-2">
            <Button asChild size="lg">
              <Link to="/dashboard">
                Open app <ArrowRight className="ml-1 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </header>

      <main>
        {/* Hero */}
        <section className="mx-auto max-w-6xl px-4 pb-20 pt-16 md:px-6 md:pt-24">
          <div className="mx-auto max-w-3xl text-center">
            <span className="inline-flex items-center gap-2 rounded-full border border-border/60 bg-card px-3 py-1 text-xs text-muted-foreground">
              <span className="h-1.5 w-1.5 rounded-full bg-success" /> New ·
              Multi-range insights
            </span>
            <h1 className="mt-6 text-4xl font-semibold tracking-tight sm:text-5xl md:text-6xl">
              Personal finance,
              <br className="hidden sm:block" /> made calm.
            </h1>
            <p className="mx-auto mt-5 max-w-xl text-base text-muted-foreground md:text-lg">
              Lumen brings every account, transaction and budget into one
              focused workspace — so you can spend less time tracking and more
              time living.
            </p>
            <div className="mt-8 flex items-center justify-center gap-3">
              <Button asChild size="lg">
                <Link to="/dashboard">
                  Get started <ArrowRight className="ml-1 h-4 w-4" />
                </Link>
              </Button>
              <Button asChild size="lg" variant="outline">
                <a href="#features">See features</a>
              </Button>
            </div>
          </div>

          {/* Preview card */}
          <div className="relative mx-auto mt-16 max-w-5xl">
            <div className="rounded-2xl border border-border/60 bg-card p-2 shadow-elevated">
              <div className="rounded-xl border border-border/60 bg-background p-6">
                <div className="flex items-center justify-between border-b border-border/60 pb-4">
                  <div>
                    <p className="text-xs text-muted-foreground">Net worth</p>
                    <p className="mt-1 text-2xl font-semibold tabular-nums">
                      $284,210
                    </p>
                  </div>
                  <span className="rounded-full bg-success/15 px-2 py-0.5 text-xs font-medium text-success">
                    +4.82%
                  </span>
                </div>
                <div className="grid gap-4 pt-5 sm:grid-cols-3">
                  {[
                    {
                      label: "Cash",
                      value: "$42,180",
                      delta: "+1.2%",
                      up: true,
                    },
                    {
                      label: "Investments",
                      value: "$201,540",
                      delta: "+5.8%",
                      up: true,
                    },
                    {
                      label: "Liabilities",
                      value: "$13,800",
                      delta: "-0.6%",
                      up: false,
                    },
                  ].map((s) => (
                    <div
                      key={s.label}
                      className="rounded-lg border border-border/60 p-4"
                    >
                      <p className="text-xs text-muted-foreground">{s.label}</p>
                      <p className="mt-1 text-lg font-semibold tabular-nums">
                        {s.value}
                      </p>
                      <p
                        className={`mt-0.5 text-xs ${s.up ? "text-success" : "text-destructive"}`}
                      >
                        {s.delta}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Features */}
        <section
          id="features"
          className="border-t border-border/60 bg-muted/30"
        >
          <div className="mx-auto max-w-6xl px-4 py-20 md:px-6">
            <div className="max-w-2xl">
              <p className="text-sm text-muted-foreground">
                Everything you need
              </p>
              <h2 className="mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
                A calmer way to manage money.
              </h2>
            </div>
            <div className="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {features.map((f) => (
                <div
                  key={f.title}
                  className="rounded-2xl border border-border/60 bg-card p-6 transition-colors hover:bg-accent/40"
                >
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-foreground/5 text-foreground">
                    <f.icon className="h-5 w-5" />
                  </div>
                  <h3 className="mt-4 text-base font-semibold">{f.title}</h3>
                  <p className="mt-1.5 text-sm text-muted-foreground">
                    {f.desc}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" className="border-t border-border/60">
          <div className="mx-auto max-w-6xl px-4 py-20 md:px-6">
            <div className="text-center">
              <p className="text-sm text-muted-foreground">Pricing</p>
              <h2 className="mt-2 text-3xl font-semibold tracking-tight md:text-4xl">
                Start free. Upgrade when you're ready.
              </h2>
            </div>
            <div className="mx-auto mt-12 grid max-w-4xl gap-4 md:grid-cols-2">
              {[
                {
                  name: "Free",
                  price: "$0",
                  desc: "For getting your finances in one place.",
                  features: [
                    "Up to 3 accounts",
                    "Budgets & transactions",
                    "Light & dark themes",
                  ],
                  cta: "Get started",
                },
                {
                  name: "Plus",
                  price: "$8",
                  desc: "For people who want the full picture.",
                  features: [
                    "Unlimited accounts",
                    "Investment tracking",
                    "Custom categories",
                    "Priority support",
                  ],
                  cta: "Try Plus",
                  highlight: true,
                },
              ].map((p) => (
                <div
                  key={p.name}
                  className={`rounded-2xl border p-6 ${p.highlight ? "border-foreground bg-card shadow-elevated" : "border-border/60 bg-card"}`}
                >
                  <div className="flex items-baseline justify-between">
                    <h3 className="text-lg font-semibold">{p.name}</h3>
                    <p className="text-2xl font-semibold tabular-nums">
                      {p.price}
                      <span className="text-sm font-normal text-muted-foreground">
                        /mo
                      </span>
                    </p>
                  </div>
                  <p className="mt-1 text-sm text-muted-foreground">{p.desc}</p>
                  <ul className="mt-5 space-y-2 text-sm">
                    {p.features.map((f) => (
                      <li
                        key={f}
                        className="flex items-center gap-2 text-muted-foreground"
                      >
                        <span className="h-1.5 w-1.5 rounded-full bg-foreground/60" />{" "}
                        {f}
                      </li>
                    ))}
                  </ul>
                  <Button
                    asChild
                    className="mt-6 w-full"
                    variant={p.highlight ? "default" : "outline"}
                  >
                    <Link to="/dashboard">{p.cta}</Link>
                  </Button>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* CTA */}
        <section className="border-t border-border/60 bg-muted/30">
          <div className="mx-auto max-w-4xl px-4 py-20 text-center md:px-6">
            <h2 className="text-3xl font-semibold tracking-tight md:text-4xl">
              Ready to see the whole picture?
            </h2>
            <p className="mx-auto mt-3 max-w-xl text-muted-foreground">
              Open the app and explore the demo dashboard — no signup required.
            </p>
            <Button asChild size="lg" className="mt-8">
              <Link to="/dashboard">
                Open dashboard <ArrowRight className="ml-1 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </section>
      </main>

      <footer id="faq" className="border-t border-border/60">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 px-4 py-8 text-sm text-muted-foreground md:flex-row md:px-6">
          <p>© {new Date().getFullYear()} Lumen. All rights reserved.</p>
          <div className="flex items-center gap-5">
            <a href="#" className="hover:text-foreground">
              Privacy
            </a>
            <a href="#" className="hover:text-foreground">
              Terms
            </a>
            <a href="#" className="hover:text-foreground">
              Contact
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
