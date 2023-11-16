
const C_RED = '#CA4538';
const C_BLUE = '#417CBF';
const C_GREEN = '#337B3D';
const C_YELLOW = '#EF8E3C';
const C_PINK = '#E94F83';
const C_ORANGE = '#EB582D';
const C_GRAY = '#7A8486';
const C_PURPLE = '#7E45C5';

const C_EMPTY =  'rgb(219 234 254)';

const GAME_WIDTH = 720;
const GAME_HEIGHT = 480;

export class GameGrid {
    constructor(x, y, picks = new Set([])) {
        this.x = x;
        this.y = y;

        this.cols = 10;
        this.rows = 8; // Total
        this.group_rows = 4; // Grouped (HEADS, TAILS)

        this.game_id = 0;
        this.picks = picks;
        this.nextGameTime = Date.now();

        this.newNum = 0;
        this.newNumTime = 0;
        this.newNumAnimationDuration = 1000;
        this.newNumAnimationWaitDuration = 500;
    }

    update(game_id, picks, nextGameTime) {
        // Get the new number
        if (this.picks.size === picks.size - 1) {
            picks.forEach((pick) => {
                if (!this.picks.has(pick)) {
                    this.newNum = pick;
                    this.newNumTime = Date.now();
                }
            });
        }

        // Update the game state
        this.game_id = game_id;
        this.nextGameTime = nextGameTime;
        this.picks = new Set([]);
        picks.forEach((pick) => { this.picks.add(pick); });
    }

    lerp(a, b, t) {
        return a + (b - a) * t;
    }

    render(ctx) {
        // Render Header
        let hdrHeight = this.renderHeader(ctx);

        // Render Heads and Tails Groupings
        let headHeight = this.renderGroup(ctx, false, hdrHeight);
        this.renderGroup(ctx, true, hdrHeight + headHeight);

        // Animate New Number
        if (this.newNumTime > 0) {
            const now = Date.now();
            const timeLeftWaiting = this.newNumTime + this.newNumAnimationWaitDuration - now;
            const timeLeft = this.newNumTime + this.newNumAnimationDuration + this.newNumAnimationWaitDuration - now;

            const row = Math.floor((this.newNum - 1) / this.cols);
            const col = (this.newNum - 1) % this.cols;

            const cellX = col * (64 + 4) + 40 + 4 + (64/2);
            const cellY = hdrHeight + row * (32 + 4) + 4 + (this.newNum > 40 ? 8 : 0) + (32/2);

            if (timeLeft > 0) {

                let radius = 100;
                let x = this.x + GAME_WIDTH/2;
                let y = this.y + hdrHeight + headHeight;
                let fontSize = 72;
                if (timeLeftWaiting < 0) {
                    radius = 100 * (timeLeft / this.newNumAnimationDuration);
                    x = this.lerp(this.x + GAME_WIDTH/2, this.x + cellX, 1 - (timeLeft / this.newNumAnimationDuration));
                    y = this.lerp(this.y + hdrHeight + headHeight, this.y + cellY, 1 - (timeLeft / this.newNumAnimationDuration)) - 5;
                    fontSize = 72 * (timeLeft / this.newNumAnimationDuration);
                }

                ctx.beginPath();
                ctx.fillStyle = this.fillColour(row);
                ctx.arc(x, y, radius, 0, 2 * Math.PI);
                ctx.fill();

                // Draw the number
                ctx.font = `bold ${fontSize}px Arial`;
                ctx.fillStyle = 'white';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillText(this.newNum, x, y);

            } else {
                this.newNumTime = 0;
            }
        }
    }

    getHeads() {
        let heads = 0;
        this.picks.forEach((pick) => {
            if (pick <= 40) {
                heads++;
            }
        });
        return heads;
    }
    getTails() {
        let tails = 0;
        this.picks.forEach((pick) => {
            if (pick > 40) {
                tails++;
            }
        });
        return tails;
    }
    getTimeLeft() {
        const now = Date.now();
        const timeLeft = this.nextGameTime - now;
        if (timeLeft < 0) {
            return "00:00";
        }

        const minutes = Math.floor(timeLeft / 60_000);
        const seconds = Math.floor((timeLeft - minutes * 60_000) / 1_000);

        return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }

    fillColour(row) {
        const colors = [
            '#CA4538',
            '#417CBF',
            '#337B3D',
            '#EF8E3C',
            '#E94F83',
            '#EB582D',
            '#7A8486',
            '#7E45C5',
        ];

        return colors[row % colors.length];
    };

    renderGroup(ctx, tails, y_offset) {
        const cell_width = 64;
        const cell_height = 32;
        const cell_padding = 4;
        const cell_radius = 2;

        const group_height = this.group_rows * (cell_height + cell_padding) - cell_padding;
        const group_padding = 8;

        // Draw Heads
        ctx.beginPath();
        ctx.fillStyle = tails ? C_BLUE : C_RED;
        ctx.roundRect(this.x, this.y + y_offset, 40, group_height, [cell_radius]);
        ctx.fill();
        
        ctx.save();
        ctx.translate(this.x + 20, this.y + group_height/2 + y_offset);
            ctx.rotate(-Math.PI/2);
            ctx.font = '18px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText( tails ? "Tails" : "Heads", 0, 0);
        ctx.restore();

        // Draw Heads cells
        ctx.save();
        ctx.translate(this.x + 40 + cell_padding, this.y + y_offset);
        for (let col = 0; col < this.cols; col++) {
            for (let row = 0; row < this.group_rows; row++) {
                const x = col * (cell_width + cell_padding);
                const y = row * (cell_height + cell_padding);
                const width = cell_width;
                const height = cell_height;

                // Check if the cell has been picked
                const pick = (col + 1) + (row + (tails ? 4 : 0)) * this.cols;
                let picked = false; 
                this.picks.forEach((p) => picked |= (p === pick));
                
                if (picked) {

                    // Check if the cell is the new number wait for animation to finish
                    if (this.newNum === pick && this.newNumTime > 0) {
                        const now = Date.now();
                        const timeLeft = this.newNumTime + this.newNumAnimationDuration + this.newNumAnimationWaitDuration - now;
                        if (timeLeft > 0) {
                            ctx.beginPath();
                            ctx.roundRect(x, y, width, height, [cell_radius]);
                            ctx.fillStyle = C_EMPTY;
                            ctx.fill();
                            continue;
                        }
                    }

                    ctx.beginPath();
                    ctx.roundRect(x, y, width, height, [cell_radius]);
                    ctx.fillStyle = this.fillColour(row  + (tails ? 4 : 0));
                    ctx.fill();

                    // Draw the number
                    ctx.font = '18px Arial';
                    ctx.fillStyle = 'white';
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'middle';
                    ctx.fillText(pick, x + width/2, y + height/2);
                } else {
                    ctx.beginPath();
                    ctx.roundRect(x, y, width, height, [cell_radius]);
                    ctx.fillStyle = C_EMPTY;
                    ctx.fill();
                }

            }
        }
        ctx.restore();

        return group_height + group_padding;
    }

    renderHeader(ctx) {
        // Draw Serving Current Game Number Box
        const primary_box_width = 64+64+4;
        const primary_box_height = 48;

        const secondary_box_width = 64;
        const secondary_box_height = 32;

        const box_radius = 2;

        ctx.save();
        ctx.translate(this.x + 720 - primary_box_width, this.y);

            ctx.beginPath();
            ctx.fillStyle = C_BLUE;
            ctx.roundRect(0, 0, primary_box_width, primary_box_height, [box_radius]);
            ctx.fill();

            // Draw the "DRAWING GAME" text
            ctx.font = 'bold 10px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText("DRAWING GAME", primary_box_width/2, 12);

            // Draw the Game Number "#" text
            ctx.font = 'bold 24px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText(`${this.game_id}`, primary_box_width/2, 34);

        ctx.restore();

        ctx.save();
        ctx.translate(this.x + 720 - 2 * primary_box_width - 4, this.y);

            ctx.beginPath();
            ctx.fillStyle = C_BLUE;
            ctx.roundRect(0, 0, primary_box_width, primary_box_height, [box_radius]);
            ctx.fill();

            // Draw the "DRAWING GAME" text
            ctx.font = 'bold 10px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText("NEXT GAME", primary_box_width/2, 12);

            // Draw the Game Number "#" text
            ctx.font = 'bold 24px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText(`${this.getTimeLeft()}`, primary_box_width/2, 34);

        ctx.restore();

        ctx.save();
        ctx.translate(this.x + 720 - 2 * primary_box_width - 2 * 4 - secondary_box_width, this.y);

            ctx.beginPath();
            ctx.fillStyle = C_BLUE;
            ctx.roundRect(0, primary_box_height - secondary_box_height, secondary_box_width, secondary_box_height, [box_radius]);
            ctx.fill();

            // Draw the "DRAWING GAME" text
            ctx.font = 'bold 12px Arial';
            ctx.fillStyle = C_BLUE;
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText("Tails", secondary_box_width/2, 10);

            // Draw the Game Number "#" text
            ctx.font = 'bold 24px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText(`${this.getTails()}`, secondary_box_width/2, 34);

        ctx.restore();

        ctx.save();
        ctx.translate(this.x + 720 - 2 * primary_box_width - 3 * 4 - 2 * secondary_box_width, this.y);

            ctx.beginPath();
            ctx.fillStyle = C_RED;
            ctx.roundRect(0, primary_box_height - secondary_box_height, secondary_box_width, secondary_box_height, [box_radius]);
            ctx.fill();

            // Draw the "DRAWING GAME" text
            ctx.font = 'bold 12px Arial';
            ctx.fillStyle = C_RED;
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText("Heads", secondary_box_width/2, 10);

            // Draw the Game Number "#" text
            ctx.font = 'bold 24px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText(`${this.getHeads()}`, secondary_box_width/2, 34);

        ctx.restore();

        return primary_box_height + 4;
    }
}