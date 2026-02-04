# Original Idea

## The Trade

I will buy Oracle call ratio spreads when the time is right, I will look to put on 3-4x ratios: sell slightly otm put to fund deep otm calls. This is what is called a vol skew trade, specifically I am expressing that I believe downside risk is overpriced and I am buying underpriced asymmetric upside risk (you can watch implied vol which is a function of price to determine relative overall value, currently slightly otm puts trade ~60% impVol @ 1y dte while deep otm calls trade ~54%).

## Skew Analysis

Another way to look at it is take delta equivalent otm put and otm call, compare impVol, ~62% put impVol 12c delta & ~53% call impVol 12c delta, current px @ $153 | call strike @ $330 | put strike @ 95.. this tells you the reversal is currently priced @ -9% which means there is put skew signaling puts are rich relative to calls, puts tend to be structurally rich relative to calls bc real money accounts hedge downside for portfolios programatically, you chart reversal ratio over time + you map nominal delta neutral px difference in this case 330-95 = 235 and track that as well, you do that across all deltas (this is just one example using a ~12c delta). You are now trading option skew otherwise known as gamma, my strategy proposed above specifically is referred to as "picking gamma". Can be 1:1 or 3-4:1 as I suggested depends on the risk profile you prefer to tune.

## Smile Analysis

You do this same approach but not delta equivalent but rather nominal strike equivalent and you are now analyzing what is called vol smile.

- Skew = directionality tells
- Smile = statistical tail tells

The shape of the smile reveals tail directionality which allows you to derive probabilistic odds using Bayesian priors

Probabilistic odds allow you to use Kelly criterion in order to determine sizing (think poker and how your hand relative to the board and odds of winning determine the % of your pot you should be willing to wager on any given hand)

## Volatility Surface

Now you do skew + smile analysis on each expiry, and now you have a 3d vol surface

## Heston Model

And now it gets fun! We now fit a Heston model

(black-scholes assumes constant vol, Heston assumes it is stochastically random which best empirically explains the mkt, why? Bc the market displays informatic unpredictability best explained by Brownian motion)

we fit the Heston model using the vol current vol skew & smile, this model now allows us to simulate forward surfaces, estimate tail probabilities, estimate forward skew, and extract risk neutral densities in the mkt

we fit/train the model iteratively daily

## Relative Value

surface tells us what the mkt currently believes in regards to price, Heston tells us exactly how the mkt believes volatility will behave… we now have a way to track historical vols and compare current pricing vol relative to the models implied vol

we have now built a options vol relative value pricing model which we could use programmatically

## Vision

The above acts as a primer, I share it as a means to help dispense the knowledge so that you both can review at your convenience, and as a way to keep me accountable, I should build this for us in the future. This is the foundational groundwork for the money printing machine, we tune risk down such that we only take highly asymmetric trades. We run this strategy on a multitude of stocks. We couple this with momentum studies since momentum is directly correlated to vol skew/smile mispricing.
